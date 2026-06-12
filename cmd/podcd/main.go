package main

import (
	"context"
	"log"

	"github.com/juli3nk/podcd/internal/controller"
	initrunner "github.com/juli3nk/podcd/internal/init"
	"github.com/juli3nk/podcd/internal/reconcile"
	"github.com/juli3nk/podcd/internal/renderer"
	"github.com/juli3nk/podcd/internal/runtime"
	"github.com/juli3nk/podcd/internal/source"
	"github.com/juli3nk/podcd/internal/state"
	"github.com/juli3nk/podcd/internal/systemd"
)

func main() {
	ctx := context.Background()

	// 1. Source (Git)
	src, err := source.NewGit(
		"git@github.com:me/repo.git",
		"main",
		"/opt/podcd/repo",
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Runtime détecté
	rt, err := runtime.DetectRuntime()
	if err != nil {
		log.Fatal(err)
	}

	// 3. Renderer (docker ou podman)
	rend, err := renderer.New(rt)
	if err != nil {
		log.Fatal(err)
	}

	// 4. State store
	store, err := state.NewFileStateStore("/var/lib/podcd/state.json")
	if err != nil {
		log.Fatal(err)
	}

	// 5. Init runner
	initRunner, err := initrunner.New(rt, store)
	if err != nil {
		log.Fatal(err)
	}

	// 6. systemd manager
	sysd := &systemd.Systemd{
		UserMode: false, // ou true pour laptop
	}

	// 7. Reconciler
	reconciler := reconcile.New(
		src,
		rend,
		sysd,
		initRunner,
	)

	// 8. Controller
	ctrl := controller.New(reconciler)

	// 9. Run
	ctrl.Run(ctx)
}
