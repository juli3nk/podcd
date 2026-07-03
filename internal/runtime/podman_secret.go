package runtime

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/juli3nk/podcd/internal/model"
)

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

func (r *PodmanRuntime) CreateSecret(secret model.Secret, data []byte, hash string) error {
	args := []string{"secret", "create"}

	labels := managedLabels(secret.Name, hash, secret.Labels)
	for k, val := range labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, val))
	}

	args = append(args, secret.Name, "-")

	cmd := exec.Command(podmanExec, args...)
	cmd.Stdin = bytes.NewReader(data)

	return cmd.Run()
}

func (r *PodmanRuntime) RemoveSecret(name string) error {
	args := []string{"secret", "remove"}

	args = append(args, name)

	return exec.Command(podmanExec, args...).Run()
}
