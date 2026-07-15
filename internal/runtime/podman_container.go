package runtime

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/juli3nk/podcd/internal/model"
)

func (r *PodmanRuntime) ListContainers(filter Labels) ([]ContainerInfo, error) {
	args := []string{"container", "ls", "-a", "--format", "{{.ID}}"}

	for k, v := range filter {
		args = append(args, "--filter", fmt.Sprintf("label=%s=%s", k, v))
	}

	out, err := exec.Command(r.binaryPath, args...).Output()
	if err != nil {
		return nil, err
	}

	var result []ContainerInfo

	for _, id := range strings.Fields(string(out)) {
		info, err := r.inspectContainer(id)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	return result, nil
}

func (r *PodmanRuntime) Run(spec model.Container, hash string) error {
	args := BuildPodmanContainerRunArgs(spec, hash)

	return exec.Command(r.binaryPath, args...).Run()
}

func (r *PodmanRuntime) inspectContainer(id string) (ContainerInfo, error) {
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

func BuildPodmanContainerRunArgs(spec model.Container, hash string) []string {
	args := []string{"container", "run"}

	if spec.Remove {
		args = append(args, "--rm")
	}

	for _, v := range spec.Volumes {
		mount := fmt.Sprintf("type=%s,src=%s,dst=%s", v.Type, v.Source, v.Target)
		if v.ReadOnly {
			mount = fmt.Sprintf("%s,ro", mount)
		}
		args = append(args, "--mount", mount)
	}

	for _, n := range spec.Networks {
		args = append(args, "--network", n.Name)
	}

	for _, d := range spec.DNS {
		args = append(args, "--dns", d)
	}

	for _, p := range spec.Ports {
		args = append(args, fmt.Sprintf("--publish %d:%d", p.HostPort, p.ContainerPort))
	}

	for k, v := range spec.Env {
		args = append(args, fmt.Sprintf("--env %s=%s", k, v))
	}

	for _, cm := range spec.ConfigMaps {
		args = append(args, "--mount", fmt.Sprintf("type=tmpfs,src=%s,dst=%s,ro", cm.Name, cm.Target))
	}

	for _, s := range spec.Secrets {
		args = append(args, "--secret", fmt.Sprintf("%s,type=mount,target=%s", s.Name, s.Target))
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

	return args
}
