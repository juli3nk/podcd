package model

type InitTask struct {
	Name string

	Volumes []VolumeRef
	Devices []string

	Networks []NetworkRef
	DNS      []string

	Env        map[string]string
	ConfigMaps []FileRef
	Secrets    []FileRef

	Image   string
	Command []string
	Args    []string

	Once bool
}

func (it InitTask) ToContainer() Container {
	return Container{
		Name: it.Name,

		Remove: true,

		Volumes: it.Volumes,
		Devices: it.Devices,

		Networks: it.Networks,
		DNS:      it.DNS,

		Env:        it.Env,
		ConfigMaps: it.ConfigMaps,
		Secrets:    it.Secrets,

		Image:   it.Image,
		Command: it.Command,
		Args:    it.Args,
	}
}
