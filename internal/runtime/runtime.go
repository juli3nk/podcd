package runtime

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/juli3nk/podcd/internal/model"
)

type RuntimeBase struct {
	binaryPath string
	storageDir string
}

type Runtime interface {
	ListNetworks(filter Labels) ([]NetworkInfo, error)
	CreateNetwork(net model.Network, hash string) error
	RemoveNetwork(name string) error

	ListVolumes(filter Labels) ([]VolumeInfo, error)
	CreateVolume(vol model.Volume, hash string) error
	RemoveVolume(name string) error

	ListConfigMaps(filter Labels) ([]ConfigMapInfo, error)
	CreateConfigMap(configMap model.ConfigMap, hash string) error
	RemoveConfigMap(name string) error

	ListSecrets(filter Labels) ([]SecretInfo, error)
	CreateSecret(secret model.Secret, data []byte, hash string) error
	RemoveSecret(name string) error

	ListContainers(filter Labels) ([]ContainerInfo, error)
	Run(spec model.Container, hash string) error
}

func New(rt Backend, storageDir string) (Runtime, error) {
	switch rt {
	case BackendDocker:
		binaryPath, err := rt.BinaryPath()
		if err != nil {
			return nil, err
		}

		return &DockerRuntime{
			RuntimeBase: RuntimeBase{
				binaryPath: binaryPath,
				storageDir: storageDir,
			},
		}, nil
	case BackendPodman:
		binaryPath, err := rt.BinaryPath()
		if err != nil {
			return nil, err
		}

		return &PodmanRuntime{
			RuntimeBase: RuntimeBase{
				binaryPath: binaryPath,
				storageDir: storageDir,
			},
		}, nil
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

func (r *RuntimeBase) secretDir() string {
	return filepath.Join(
		r.storageDir,
		"secrets",
	)
}

func (r *RuntimeBase) secretPath(name string) string {
	return filepath.Join(
		r.secretDir(),
		name,
	)
}

func (r *RuntimeBase) configMapDir() string {
	return filepath.Join(
		r.storageDir,
		"configmaps",
	)
}

func (r *RuntimeBase) configMapPath(name string) string {
	return filepath.Join(
		r.configMapDir(),
		name,
	)
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
	result[string(LabelResourceHash)] = hash

	return result
}
