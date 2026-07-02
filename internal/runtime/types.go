package runtime

type Backend string

const (
	BackendDocker Backend = "docker"
	BackendPodman Backend = "podman"
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
	LabelManaged  LabelType = "podcd.io/managed"
	LabelName     LabelType = "podcd.io/name"
	LabelSpecHash LabelType = "podcd.io/spec-hash"
)

type Labels map[LabelType]string

type ResourceInfo struct {
	Name   string
	Labels Labels
}

type ContainerInfo = ResourceInfo
type NetworkInfo = ResourceInfo
type SecretInfo = ResourceInfo
type VolumeInfo = ResourceInfo

type inspectMeta struct {
	Name   string
	Labels Labels
}
