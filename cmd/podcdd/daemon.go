package main

import (
	"context"
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

	"golang.org/x/sync/errgroup"
)

func runDaemon(cmd *cobra.Command, args []string) {
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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Config
	userMode := config.IsUserMode()
	paths := config.DefaultPaths(userMode)

	cfg, err := config.LoadConfig(paths.Config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if level, err := zerolog.ParseLevel(cfg.Server.LogLevel); err == nil {
		zerolog.SetGlobalLevel(level)
	}

	if err := config.EnsureDirectories(paths); err != nil {
		log.Fatal().Err(err).Msg("failed to create directories")
	}

	// Identity
	idm := identity.New(paths.IdentityDir)

	if err := idm.EnsureSSHKey(); err != nil {
		log.Fatal().Err(err).Msg("failed to ensure ssh key")
	}

	if err := idm.EnsureAgeKey(); err != nil {
		log.Fatal().Err(err).Msg("failed to ensure age key")
	}

	// Source
	src, err := source.NewGit(
		cfg.Source.RepoURL,
		cfg.Source.Version,
		paths.Workspace,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("source init")
	}

	// Secrets
	decrypter := secret.NewSOPSDecrypter(
		filepath.Join(paths.IdentityDir, "age", "keys.txt"),
	)

	// Runtime
	runtimeBackend, err := runtime.DetectRuntime()
	if err != nil {
		log.Fatal().Err(err).Msg("detect runtime")
	}

	rt, err := runtime.New(
		runtimeBackend,
		paths.RuntimeStorageDir,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("runtime init")
	}

	// Renderer
	rend, err := renderer.New(runtimeBackend)
	if err != nil {
		log.Fatal().Err(err).Msg("renderer init")
	}

	// Systemd
	sysd, err := systemd.New(userMode)
	if err != nil {
		log.Fatal().Err(err).Msg("systemd init")
	}

	// State
	store, err := state.NewFileStateStore(paths.State)
	if err != nil {
		log.Fatal().Err(err).Msg("state init")
	}

	// Init runner
	initRunner, err := initrunner.New(
		runtimeBackend,
		paths.RuntimeStorageDir,
		store,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("init runner")
	}

	// Reconciler
	reconciler := reconcile.New(
		src,
		decrypter,
		rend,
		sysd,
		rt,
		initRunner,
	)

	// Controller
	ctrl := controller.New(
		reconciler,
		cfg.Controller.Interval,
	)

	// IPC
	h := handler.New(
		runtimeBackend,
		reconciler,
		idm,
	)

	server := ipc.NewServer(
		paths.Socket,
		h.Handle,
	)

	log.Info().
		Str("runtime", cfg.Runtime.Type).
		Bool("userMode", userMode).
		Msgf("starting %s service", appName)

	g, ctx := errgroup.WithContext(ctx)

	// Controller
	g.Go(func() error {
		return ctrl.Run(ctx)
	})

	// IPC
	g.Go(func() error {
		return server.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Error().
			Err(err).
			Msg("service terminated unexpectedly")
	} else {
		log.Info().
			Msg("shutdown complete")
	}

	log.Info().Msgf("%s service stopped", appName)
}
