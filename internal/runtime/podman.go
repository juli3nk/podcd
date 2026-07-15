package runtime

import (
	"os/exec"
)

type PodmanRuntime struct {
	RuntimeBase
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
