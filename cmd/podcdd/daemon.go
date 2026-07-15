package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/juli3nk/podcd/internal/config"
	"github.com/juli3nk/podcd/internal/controller"
	"github.com/juli3nk/podcd/internal/handler"
	"github.com/juli3nk/podcd/internal/identity"
	initrunner "github.com/juli3nk/podcd/internal/init"
	"github.com/juli3nk/podcd/internal/ipc"
	"github.com/juli3nk/podcd/internal/reconcile"
	"github.com/juli3nk/podcd/internal/renderer"
	"github.com/juli3nk/podcd/internal/runtime"
	"github.com/juli3nk/podcd/internal/secret"
	"github.com/juli3nk/podcd/internal/source"
	"github.com/juli3nk/podcd/internal/state"
	"github.com/juli3nk/podcd/internal/systemd"
)

func runDaemon(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Logger
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	log := zerolog.New(output).With().
		Timestamp().
		Str("service", appName).
		Str("version", "1.0.0").
		Logger()

	// Config
	userMode := config.IsUserMode()
	paths := config.DefaultPaths(userMode)

	cfg, err := config.LoadConfig(paths.Config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Set log level
	level, err := zerolog.ParseLevel(cfg.Server.LogLevel)
	if err == nil {
		zerolog.SetGlobalLevel(level)
	}

	if err := config.EnsureDirectories(paths); err != nil {
		log.Fatal().Err(err).Msg("failed to create directories")
	}

	idm := identity.New(paths.IdentityDir)

	if err := idm.EnsureSSHKey(); err != nil {
		log.Fatal().Err(err).Msg("failed to ensure ssh key")
	}

	if err := idm.EnsureAgeKey(); err != nil {
		log.Fatal().Err(err).Msg("failed to ensure age key")
	}

	// 1. Source (Git)
	src, err := source.NewGit(
		cfg.Source.RepoURL,
		cfg.Source.Version,
		paths.Workspace)
	if err != nil {
		log.Fatal().Err(err).Msg("source init")
	}

	// Decrypter
	decrypter := secret.NewSOPSDecrypter(filepath.Join(paths.IdentityDir, "age", "keys.txt"))

	// 2. Runtime
	runtimeBackend, err := runtime.DetectRuntime()
	if err != nil {
		log.Fatal().Err(err).Msg("detect runtime")
	}

	rt, err := runtime.New(runtimeBackend, paths.RuntimeStorageDir)
	if err != nil {
		log.Fatal().Err(err).Msg("runtime init")
	}

	// 3. Renderer
	rend, err := renderer.New(runtimeBackend)
	if err != nil {
		log.Fatal().Err(err).Msg("renderer init")
	}

	// 4. systemd manager
	// sysd := &systemd.Systemd{
	// 	UserMode: userMode,
	// }
	sysd, err := systemd.New(userMode)
	if err != nil {
		log.Fatal().Err(err).Msg("systemd init")
	}

	// 6. State
	store, err := state.NewFileStateStore(paths.State)
	if err != nil {
		log.Fatal().Err(err).Msg("state init")
	}

	// 7. Init runner
	initRunner, err := initrunner.New(runtimeBackend, paths.RuntimeStorageDir, store)
	if err != nil {
		log.Fatal().Err(err).Msg("init runner")
	}

	// 8. Reconciler
	reconciler := reconcile.New(
		src,
		decrypter,
		rend,
		sysd,
		rt,
		initRunner,
	)

	// 9. Controller
	ctrl := controller.New(reconciler, cfg.Controller.Interval)

	log.Info().
		Str("runtime", cfg.Runtime.Type).
		Bool("userMode", userMode).
		Msgf("Starting %s service", appName)

	// 10. Run
	done := make(chan error, 1)

	go func() {
		done <- ctrl.Run(ctx)
	}()

	h := handler.New(
		runtimeBackend,
		reconciler,
		idm,
	)

	server := ipc.NewServer(
		paths.Socket,
		h.Handle,
	)

	go func() {
		if err := server.Start(); err != nil {
			log.Error().Err(err).Msg("ipc server stopped")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Info().Str("signal", sig.String()).Msg("shutdown requested")
	case err := <-done:
		fmt.Printf("%+v", err)
		log.Error().Err(err).Msg("controller stopped unexpectedly")
	}

	cancel()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		log.Warn().Msg("force shutdown after timeout")
	}

	log.Info().Msgf("%s service stopped", appName)
}
