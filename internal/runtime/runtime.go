package runtime

import (
	"fmt"
	"os/exec"

	"github.com/juli3nk/podcd/internal/model"
)

type RunSpec struct {
	Image string

	Command []string
	Args    []string

	Env map[string]string

	Volumes  []model.VolumeRef
	Networks []model.NetworkRef

	Remove bool
	Name   string
}

type Runtime interface {
	Run(spec RunSpec) error

	CreateVolume(vol model.Volume) error
	CreateNetwork(net model.Network) error
}

func New(rt model.RuntimeType) (Runtime, error) {
	switch rt {
	case model.RuntimeDocker:
		return &DockerRuntime{}, nil
	case model.RuntimePodman:
		return &PodmanRuntime{}, nil
	default:
		return nil, fmt.Errorf("unsupported runtime: %s", rt)
	}
}

func DetectRuntime() (model.RuntimeType, error) {
	if _, err := exec.LookPath("podman"); err == nil {
		return model.RuntimePodman, nil
	}

	if _, err := exec.LookPath("docker"); err == nil {
		return model.RuntimeDocker, nil
	}

	return "", fmt.Errorf("no supported runtime found")
}

func FromContainer(spec model.ContainerSpec) RunSpec {
	return RunSpec{
		Image:    spec.Image,
		Command:  spec.Command,
		Env:      spec.Env,
		Volumes:  spec.Volumes,
		Networks: spec.Networks,
	}
}
