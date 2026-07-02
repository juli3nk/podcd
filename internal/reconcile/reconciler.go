package reconcile

import (
	initrunner "github.com/juli3nk/podcd/internal/init"
	"github.com/juli3nk/podcd/internal/renderer"
	"github.com/juli3nk/podcd/internal/runtime"
	"github.com/juli3nk/podcd/internal/source"
	"github.com/juli3nk/podcd/internal/systemd"
)

type ResourceType string

const (
	ResourceContainer ResourceType = "container"
	ResourceNetwork   ResourceType = "network"
	ResourceSecret    ResourceType = "secret"
	ResourceVolume    ResourceType = "volume"
)

type ContainerReconciler struct {
	systemd  systemd.Manager
	renderer renderer.Renderer
	runtime  runtime.Runtime
	init     initrunner.InitRunner
}

type NetworkReconciler struct {
	runtime runtime.Runtime
}

type SecretReconciler struct {
	runtime runtime.Runtime
}

type VolumeReconciler struct {
	runtime runtime.Runtime
}

type Reconciler struct {
	source source.Source

	containers *ContainerReconciler
	networks   *NetworkReconciler
	secrets    *SecretReconciler
	volumes    *VolumeReconciler
}

func New(
	src source.Source,
	renderer renderer.Renderer,
	systemd systemd.Manager,
	runtime runtime.Runtime,
	init initrunner.InitRunner,
) *Reconciler {
	return &Reconciler{
		source: src,
		containers: &ContainerReconciler{
			systemd:  systemd,
			renderer: renderer,
			runtime:  runtime,
			init:     init,
		},
		networks: &NetworkReconciler{
			runtime: runtime,
		},
		secrets: &SecretReconciler{
			runtime: runtime,
		},
		volumes: &VolumeReconciler{
			runtime: runtime,
		},
	}
}
