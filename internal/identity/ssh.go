package identity

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (m *IdentityManager) EnsureSSHKey() error {
	path := filepath.Join(
		m.basePath,
		"ssh",
		"id_ed25519",
	)

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	cmd := exec.Command(
		"ssh-keygen",
		"-t", "ed25519",
		"-f", path,
		"-N", "",
	)

	return cmd.Run()
}

func (m *IdentityManager) SSHPublicKey() (string, error) {
	data, err := os.ReadFile(
		filepath.Join(m.basePath, "ssh", "id_ed25519.pub"),
	)

	return strings.TrimSpace(string(data)), err
}
