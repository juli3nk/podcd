package model

type VolumeRef struct {
	Type     string
	Source   string
	Target   string
	ReadOnly bool
}

type NetworkRef struct {
	Name    string
	IP      string
	Aliases []string
}

type PortSpec struct {
	HostPort      int
	ContainerPort int
	Protocol      string
}

type FileRef struct {
	Name   string
	Target string
}

type ContainerStatus struct {
	LastAppliedHash string
	Ready           bool
}

type Container struct {
	Name   string
	Labels map[string]string

	Init []InitTask

	Remove bool

	Volumes []VolumeRef
	Devices []string

	Networks []NetworkRef
	DNS      []string
	Ports    []PortSpec

	Env        map[string]string
	ConfigMaps []FileRef
	Secrets    []FileRef

	AddCapabilities  []string
	DropCapabilities []string

	Image   string
	Command []string
	Args    []string

	Status ContainerStatus
}
