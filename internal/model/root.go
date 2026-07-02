package model

type RootSpec struct {
	Networks   []Network
	Volumes    []Volume
	Secrets    []Secret
	Containers []Container
}
