package reconcile

import (
	initrunner "github.com/juli3nk/podcd/internal/init"
	"github.com/juli3nk/podcd/internal/renderer"
	"github.com/juli3nk/podcd/internal/source"
	"github.com/juli3nk/podcd/internal/systemd"
)

type Reconciler struct {
	source   source.Source
	renderer renderer.Renderer
	systemd  systemd.Manager
	init     initrunner.InitRunner
}

func New(
	src source.Source,
	renderer renderer.Renderer,
	systemd systemd.Manager,
	init initrunner.InitRunner,
) *Reconciler {
	return &Reconciler{
		source:   src,
		renderer: renderer,
		systemd:  systemd,
		init:     init,
	}
}
