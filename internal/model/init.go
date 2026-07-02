package model

type InitTask struct {
	Name string

	Volumes  []VolumeRef
	Devices  []string
	Networks []NetworkRef
	DNS      []string
	Secrets  []SecretRef
	Env      map[string]string

	Image   string
	Command []string

	Once bool
}

func (it InitTask) ToContainer() ContainerSpec {
	return ContainerSpec{
		Remove:   true,
		Volumes:  it.Volumes,
		Devices:  it.Devices,
		Networks: it.Networks,
		DNS:      it.DNS,
		Env:      it.Env,
		Image:    it.Image,
		Command:  it.Command,
	}
}
