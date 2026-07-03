package runtime

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/juli3nk/podcd/internal/model"
)

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

func (d *PodmanRuntime) RemoveVolume(name string) error {
	args := []string{"volume", "remove"}

	args = append(args, name)

	return exec.Command(podmanExec, args...).Run()
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
