package runtime

import (
	"os/exec"
)

type DockerRuntime struct {
	RuntimeBase
}

func dockerInspect(id string, args ...string) ([]byte, error) {
	cmdArgs := append(args, "inspect", id)

	cmd := exec.Command(dockerExec, cmdArgs...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}
