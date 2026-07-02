package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/juli3nk/podcd/internal/model"
)

type PodmanRuntime struct{}

var podmanExec string = "podman"

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

func (d *PodmanRuntime) ListNetworks(filter Labels) ([]NetworkInfo, error) {
	var result []NetworkInfo

	args := []string{"network", "ls", "--format", "{{.ID}}"}

	for k, v := range filter {
		args = append(args, "--filter", fmt.Sprintf("label=%s=%s", k, v))
	}

	out, err := exec.Command(podmanExec, args...).Output()
	if err != nil {
		return nil, err
	}

	ids := strings.Fields(string(out))

	for _, id := range ids {
		info, err := d.inspectNetwork(id)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	return result, nil
}

func (d *PodmanRuntime) ListVolumes(filter Labels) ([]VolumeInfo, error) {
	var result []VolumeInfo

	args := []string{"volume", "ls", "--format", "{{.Name}}"}

	for k, v := range filter {
		args = append(args, "--filter", fmt.Sprintf("label=%s=%s", k, v))
	}

	out, err := exec.Command(podmanExec, args...).Output()
	if err != nil {
		return nil, err
	}

	names := strings.Fields(string(out))

	for _, name := range names {
		info, err := d.inspectVolume(name)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	return result, nil
}

func (d *PodmanRuntime) ListSecrets(filter Labels) ([]SecretInfo, error) {
	var result []SecretInfo

	args := []string{"secret", "ls", "--format", "{{.Name}}"}

	for k, v := range filter {
		args = append(args, "--filter", fmt.Sprintf("label=%s=%s", k, v))
	}

	out, err := exec.Command(podmanExec, args...).Output()
	if err != nil {
		return nil, err
	}

	names := strings.Fields(string(out))

	for _, name := range names {
		info, err := d.inspectVolume(name)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	return result, nil
}

func (d *PodmanRuntime) CreateNetwork(net model.Network, hash string) error {
	args := []string{"network", "create"}

	if net.Driver != "" {
		args = append(args, "--driver", net.Driver)
	}

	for _, ipam := range net.IPAM.Config {
		if ipam.Subnet != "" {
			args = append(args, "--subnet", ipam.Subnet)
		}
		if ipam.Gateway != "" {
			args = append(args, "--gateway", ipam.Gateway)
		}
		if ipam.IPRange != "" {
			args = append(args, "--ip-range", ipam.IPRange)
		}
	}

	if net.Internal {
		args = append(args, "--internal")
	}

	if net.Attachable {
		args = append(args, "--attachable")
	}

	labels := managedLabels(net.Name, hash, net.Labels)
	for k, val := range labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, val))
	}

	args = append(args, net.Name)

	return exec.Command(podmanExec, args...).Run()
}

func (d *PodmanRuntime) CreateVolume(vol model.Volume, hash string) error {
	args := []string{"volume", "create"}

	if vol.Driver != "" {
		args = append(args, "--driver", vol.Driver)
	}

	for k, val := range vol.Options {
		args = append(args, "--option", fmt.Sprintf("%s=%s", k, val))
	}

	labels := managedLabels(vol.Name, hash, vol.Labels)
	for k, val := range labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, val))
	}

	args = append(args, vol.Name)

	return exec.Command(podmanExec, args...).Run()
}

func (d *PodmanRuntime) RemoveNetwork(name string) error {
	args := []string{"network", "remove"}

	args = append(args, name)

	return exec.Command(podmanExec, args...).Run()
}

func (d *PodmanRuntime) RemoveVolume(name string) error {
	args := []string{"volume", "remove"}

	args = append(args, name)

	return exec.Command(podmanExec, args...).Run()
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

func podmanInspect(id string, args ...string) ([]byte, error) {
	cmdArgs := append(args, "inspect", id)

	cmd := exec.Command(podmanExec, cmdArgs...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
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

func (d *PodmanRuntime) inspectNetwork(id string) (NetworkInfo, error) {
	var info NetworkInfo

	out, err := dockerInspect(id, "network")
	if err != nil {
		return info, err
	}

	var data []inspectMeta

	if err := json.Unmarshal(out, &data); err != nil {
		return info, err
	}

	if len(data) == 0 {
		return info, fmt.Errorf("network not found: %s", id)
	}

	n := data[0]

	info.Name = n.Name
	info.Labels = safeLabels(n.Labels)

	return info, nil
}

func (d *PodmanRuntime) inspectVolume(name string) (VolumeInfo, error) {
	var info VolumeInfo

	out, err := dockerInspect(name, "volume")
	if err != nil {
		return info, err
	}

	var data []inspectMeta

	if err := json.Unmarshal(out, &data); err != nil {
		return info, err
	}

	if len(data) == 0 {
		return info, fmt.Errorf("volume not found: %s", name)
	}

	v := data[0]

	info.Name = v.Name
	info.Labels = safeLabels(v.Labels)

	return info, nil
}

func (r *PodmanRuntime) CreateSecret(secret model.Secret, data []byte) error {
	cmd := exec.Command(
		"podman",
		"secret",
		"create",
		secret.Name,
		"-",
	)

	cmd.Stdin = bytes.NewReader(data)

	return cmd.Run()
}

func (r *PodmanRuntime) RemoveSecret(name string) error {
	cmd := exec.Command(
		"podman",
		"secret",
		"rm",
		name,
	)

	return cmd.Run()
}
