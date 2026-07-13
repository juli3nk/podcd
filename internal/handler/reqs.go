package handler

import (
	"os/exec"
	"strings"

	"github.com/juli3nk/podcd/internal/ipc"
)

type BinaryInfo struct {
	Path    string `json:"path,omitempty"`
	Version string `json:"version,omitempty"`
	Error   string `json:"error,omitempty"`
}

type RuntimeInfo struct {
	Detected string `json:"detected"`

	BinaryInfo
}

type Requirements struct {
	Runtime  RuntimeInfo           `json:"runtime"`
	Binaries map[string]BinaryInfo `json:"binaries"`
}

func (h *Handler) reqs(req ipc.Request) ipc.Response {
	binaries := map[string]string{
		"git":        "--version",
		"sops":       "--version",
		"age-keygen": "",
		"ssh-keygen": "",
		"systemctl":  "--version",
	}

	resultBinaries := make(map[string]BinaryInfo)

	for name, versionArg := range binaries {
		resultBinaries[name] = binaryInfo(name, versionArg)
	}

	runtimeBinary := string(h.runtime)

	runtimeInfo := RuntimeInfo{}

	if runtimeBinary != "" {
		runtimeInfo = RuntimeInfo{
			Detected:   runtimeBinary,
			BinaryInfo: binaryInfo(runtimeBinary, "--version"),
		}
	}

	return ipc.Response{
		Success: true,
		Data: Requirements{
			Runtime:  runtimeInfo,
			Binaries: resultBinaries,
		},
	}
}

func binaryInfo(name, versionArg string) BinaryInfo {
	info := BinaryInfo{}

	path, err := exec.LookPath(name)
	if err != nil {
		info.Error = err.Error()

		return info
	}

	info.Path = path

	if versionArg != "" {
		info.Version = getVersion(name, versionArg)
	}

	return info
}

func getVersion(binary string, args ...string) string {
	cmd := exec.Command(binary, args...)

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	line := strings.SplitN(string(output), "\n", 2)[0]

	return strings.TrimSpace(line)
}
