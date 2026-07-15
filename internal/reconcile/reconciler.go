package reconcile

import (
	"github.com/juli3nk/podcd/internal/generator"
	initrunner "github.com/juli3nk/podcd/internal/init"
	"github.com/juli3nk/podcd/internal/renderer"
	"github.com/juli3nk/podcd/internal/runtime"
	"github.com/juli3nk/podcd/internal/secret"
	"github.com/juli3nk/podcd/internal/source"
	"github.com/juli3nk/podcd/internal/systemd"
)

type ResourceType string

const (
	ConfigMapContainer ResourceType = "configMap"
	ResourceContainer  ResourceType = "container"
	ResourceNetwork    ResourceType = "network"
	ResourceSecret     ResourceType = "secret"
	ResourceVolume     ResourceType = "volume"
)

type ConfigMapReconciler struct {
	runtime runtime.Runtime
}

type ContainerReconciler struct {
	runtime  runtime.Runtime
	init     initrunner.InitRunner
	systemd  systemd.Manager
	renderer renderer.Renderer
}

type NetworkReconciler struct {
	runtime runtime.Runtime
}

type SecretReconciler struct {
	runtime   runtime.Runtime
	decrypter secret.Decrypter
	generator generator.Generator
}

type VolumeReconciler struct {
	runtime runtime.Runtime
}

type Reconciler struct {
	source source.Source

	configMaps *ConfigMapReconciler
	containers *ContainerReconciler
	networks   *NetworkReconciler
	secrets    *SecretReconciler
	volumes    *VolumeReconciler
}

func New(
	src source.Source,
	decrypter secret.Decrypter,
	renderer renderer.Renderer,
	systemd systemd.Manager,
	runtime runtime.Runtime,
	init initrunner.InitRunner,
) *Reconciler {
	return &Reconciler{
		source: src,
		configMaps: &ConfigMapReconciler{
			runtime: runtime,
		},
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
			decrypter: decrypter,
			runtime:   runtime,
		},
		volumes: &VolumeReconciler{
			runtime: runtime,
		},
	}
}
