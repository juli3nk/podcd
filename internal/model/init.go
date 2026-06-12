package model

type InitTask struct {
	Name string

	Image   string
	Command []string

	Env      map[string]string
	Volumes  []VolumeRef
	Networks []NetworkRef
	Secrets  []SecretRef

	Once bool
}
