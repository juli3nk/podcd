package runtime

import (
	"fmt"
	"os/exec"

	"github.com/juli3nk/podcd/internal/model"
)

type DockerRuntime struct{}

func (d *DockerRuntime) Run(spec RunSpec) error {
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

	return exec.Command("docker", args...).Run()
}

func (d *DockerRuntime) CreateVolume(v model.Volume) error {
	args := []string{"volume", "create"}

	if v.Driver != "" {
		args = append(args, "--driver", v.Driver)
	}

	for k, val := range v.Options {
		args = append(args, "--option", fmt.Sprintf("%s=%s", k, val))
	}

	for k, val := range v.Labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, val))
	}

	args = append(args, v.Name)

	return exec.Command("docker", args...).Run()
}

func (d *DockerRuntime) CreateNetwork(n model.Network) error {
	args := []string{"network", "create"}

	if n.Driver != "" {
		args = append(args, "--driver", n.Driver)
	}

	for _, ipam := range n.IPAM {
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

	if n.Internal {
		args = append(args, "--internal")
	}

	if n.Attachable {
		args = append(args, "--attachable")
	}

	for k, val := range n.Labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, val))
	}

	args = append(args, n.Name)

	return exec.Command("docker", args...).Run()
}
