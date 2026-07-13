package runtime

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/juli3nk/podcd/internal/model"
)

func (r *DockerRuntime) ListNetworks(filter Labels) ([]NetworkInfo, error) {
	var result []NetworkInfo

	args := []string{"network", "ls", "--format", "{{.ID}}"}

	for k, v := range filter {
		args = append(args, "--filter", fmt.Sprintf("label=%s=%s", k, v))
	}

	out, err := exec.Command(r.binaryPath, args...).Output()
	if err != nil {
		return nil, err
	}

	ids := strings.Fields(string(out))

	for _, id := range ids {
		info, err := r.inspectNetwork(id)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}

	return result, nil
}

func (r *DockerRuntime) CreateNetwork(net model.Network, hash string) error {
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

	return exec.Command(r.binaryPath, args...).Run()
}

func (r *DockerRuntime) RemoveNetwork(name string) error {
	args := []string{"network", "remove"}

	args = append(args, name)

	return exec.Command(r.binaryPath, args...).Run()
}

func (r *DockerRuntime) inspectNetwork(id string) (NetworkInfo, error) {
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
