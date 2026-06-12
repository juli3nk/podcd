package state

type State struct {
	InitTasks map[string]bool `json:"initTasks"`

	Containers map[string]string // hash
	Volumes    map[string]string // hash
	Networks   map[string]string // hash
	Secrets    map[string]string // hash
}
