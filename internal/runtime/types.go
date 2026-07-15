package runtime

type Backend string

const (
	BackendDocker Backend = "docker"
	BackendPodman Backend = "podman"
)

var (
	dockerExec string = "docker"
	podmanExec string = "podman"
)

func (rt Backend) IsValid() bool {
	switch rt {
	case BackendDocker, BackendPodman:
		return true
	default:
		return false
	}
}

type LabelType string

const (
	LabelManaged      LabelType = "podcd.io/managed"
	LabelName         LabelType = "podcd.io/name"
	LabelResourceHash LabelType = "podcd.io/resource-hash"
)

type Labels map[LabelType]string

type ResourceInfo struct {
	Name   string
	Labels Labels
}

type ConfigMapInfo = ResourceInfo
type ContainerInfo = ResourceInfo
type NetworkInfo = ResourceInfo
type SecretInfo = ResourceInfo
type VolumeInfo = ResourceInfo

type inspectMeta struct {
	Name   string
	Labels Labels
}
