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

type SecretRef struct {
	Name string
}

type ContainerSpec struct {
	Remove bool

	Volumes []VolumeRef
	Devices []string

	Networks []NetworkRef
	DNS      []string
	Ports    []PortSpec

	Env     map[string]string
	Secrets []SecretRef

	AddCapabilities  []string
	DropCapabilities []string

	Labels map[string]string

	Name string

	Image   string
	Command []string

	Metadata Metadata
}

type ContainerStatus struct {
	LastAppliedHash string
	Ready           bool
}

type Container struct {
	Init []InitTask

	Spec   ContainerSpec
	Status ContainerStatus
}
