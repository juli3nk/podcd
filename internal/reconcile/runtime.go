package reconcile

import "github.com/juli3nk/podcd/internal/runtime"

type RuntimeObject struct {
	Name string
	Hash string
}

type RuntimeState struct {
	Containers map[string]RuntimeObject
	Networks   map[string]RuntimeObject
	Secrets    map[string]RuntimeObject
	Volumes    map[string]RuntimeObject
}

func (r *Reconciler) discoverFromRuntime() (RuntimeState, error) {
	var state RuntimeState

	state.Containers = make(map[string]RuntimeObject)
	state.Networks = make(map[string]RuntimeObject)
	state.Secrets = make(map[string]RuntimeObject)
	state.Volumes = make(map[string]RuntimeObject)

	// Containers
	containers, err := r.containers.runtime.ListContainers(runtime.Labels{runtime.LabelManaged: "true"})
	if err != nil {
		return state, err
	}
	for _, c := range containers {
		name := c.Labels[runtime.LabelName]
		state.Containers[name] = RuntimeObject{
			Name: name,
			Hash: c.Labels[runtime.LabelSpecHash],
		}
	}

	// Networks
	networks, err := r.containers.runtime.ListNetworks(runtime.Labels{runtime.LabelManaged: "true"})
	if err != nil {
		return state, err
	}
	for _, n := range networks {
		name := n.Name
		state.Networks[name] = RuntimeObject{
			Name: name,
			Hash: n.Labels[runtime.LabelSpecHash],
		}
	}

	// Secrets
	secrets, err := r.containers.runtime.ListSecrets(runtime.Labels{runtime.LabelManaged: "true"})
	if err != nil {
		return state, err
	}
	for _, n := range secrets {
		name := n.Name
		state.Secrets[name] = RuntimeObject{
			Name: name,
			Hash: n.Labels[runtime.LabelSpecHash],
		}
	}

	// Volumes
	volumes, err := r.containers.runtime.ListVolumes(runtime.Labels{runtime.LabelManaged: "true"})
	if err != nil {
		return state, err
	}
	for _, v := range volumes {
		name := v.Name
		state.Volumes[name] = RuntimeObject{
			Name: name,
			Hash: v.Labels[runtime.LabelSpecHash],
		}
	}

	return state, nil
}
