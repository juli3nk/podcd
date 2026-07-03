package runtime

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func (d *PodmanRuntime) ListContainers(filter Labels) ([]ContainerInfo, error) {
	args := []string{"container", "ls", "-a", "--format", "{{.ID}}"}

	for k, v := range filter {
		args = append(args, "--filter", fmt.Sprintf("label=%s=%s", k, v))
	}

	out, err := exec.Command(podmanExec, args...).Output()
	if err != nil {
		return nil, err
	}

	var result []ContainerInfo

	for _, id := range strings.Fields(string(out)) {
		info, err := d.inspectContainer(id)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	return result, nil
}

func (d *PodmanRuntime) Run(spec RunSpec, hash string) error {
	args := []string{"run"}

	if spec.Remove {
		args = append(args, "--rm")
	}

	for _, v := range spec.Volumes {
		mount := fmt.Sprintf("%s:%s", v.Source, v.Target)
		args = append(args, "-v", mount)
	}

	for k, v := range spec.Env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	labels := managedLabels(spec.Name, hash, spec.Labels)
	for k, val := range labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, val))
	}

	if spec.Name != "" {
		args = append(args, spec.Name)
	}

	args = append(args, spec.Image)

	if len(spec.Command) > 0 {
		args = append(args, spec.Command...)
	}

	if len(spec.Args) > 0 {
		args = append(args, spec.Args...)
	}

	return exec.Command(podmanExec, args...).Run()
}

func (d *PodmanRuntime) inspectContainer(id string) (ContainerInfo, error) {
	var info ContainerInfo

	out, err := podmanInspect(id, "container")
	if err != nil {
		return info, err
	}

	var data []struct {
		Name   string
		Config struct {
			Labels Labels
		}
	}

	if err := json.Unmarshal(out, &data); err != nil {
		return info, err
	}

	if len(data) == 0 {
		return info, fmt.Errorf("container not found: %s", id)
	}

	c := data[0]

	info.Name = strings.TrimPrefix(c.Name, "/")
	info.Labels = safeLabels(c.Config.Labels)

	return info, nil
}
