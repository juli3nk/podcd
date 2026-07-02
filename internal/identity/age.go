package identity

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (m *IdentityManager) EnsureAgeKey() error {
	path := filepath.Join(
		m.basePath,
		"age",
		"keys.txt",
	)

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	cmd := exec.Command(
		"age-keygen",
		"-o",
		path,
	)

	return cmd.Run()
}

func (m *IdentityManager) AgePublicKey() (string, error) {
	data, err := os.ReadFile(
		filepath.Join(m.basePath, "age", "keys.txt"),
	)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "# public key: ") {
			return strings.TrimPrefix(
				line,
				"# public key: ",
			), nil
		}
	}

	return "", fmt.Errorf("public key not found")
}
