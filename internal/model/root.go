package model

type RootSpec struct {
	ConfigMaps []ConfigMap
	Networks   []Network
	Volumes    []Volume
	Secrets    []Secret
	Containers []Container
}
