package runtime

import (
	"fmt"
	"os/exec"

	"github.com/juli3nk/podcd/internal/model"
)

type RunSpec struct {
	Remove bool

	Volumes []model.VolumeRef
	Devices []string

	Networks []model.NetworkRef
	DNS      []string
	Ports    []model.PortSpec

	Env map[string]string

	Labels map[string]string

	Name string

	Image string

	Command []string
	Args    []string
}

type Runtime interface {
	ListNetworks(filter Labels) ([]NetworkInfo, error)
	CreateNetwork(net model.Network, hash string) error
	RemoveNetwork(name string) error

	ListVolumes(filter Labels) ([]VolumeInfo, error)
	CreateVolume(vol model.Volume, hash string) error
	RemoveVolume(name string) error

	ListSecrets(filter Labels) ([]SecretInfo, error)
	CreateSecret(secret model.Secret, data []byte, hash string) error
	RemoveSecret(name string) error

	ListContainers(filter Labels) ([]ContainerInfo, error)
	Run(spec RunSpec, hash string) error
}

func New(rt Backend) (Runtime, error) {
	switch rt {
	case BackendDocker:
		return &DockerRuntime{
			basePath: "",
		}, nil
	case BackendPodman:
		return &PodmanRuntime{}, nil
	default:
		return nil, fmt.Errorf("unsupported runtime: %s", rt)
	}
}

func DetectRuntime() (Backend, error) {
	if _, err := exec.LookPath(dockerExec); err == nil {
		return BackendDocker, nil
	}

	if _, err := exec.LookPath(podmanExec); err == nil {
		return BackendPodman, nil
	}

	return "", fmt.Errorf("no supported runtime found")
}

func (b Backend) BinaryPath() (string, error) {
	return exec.LookPath(string(b))
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

func safeLabels(labels Labels) Labels {
	if labels == nil {
		return Labels{}
	}
	return labels
}

func managedLabels(name, hash string, labels map[string]string) map[string]string {
	result := make(map[string]string)

	for k, v := range labels {
		result[k] = v
	}

	result[string(LabelManaged)] = "true"
	result[string(LabelName)] = name
	result[string(LabelSpecHash)] = hash

	return result
}
