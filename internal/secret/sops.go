package secret

import (
	"os/exec"
)

type SOPSDecrypter struct {
	ageKeyFile string
}

func NewSOPSDecrypter(ageKeyFile string) *SOPSDecrypter {
	return &SOPSDecrypter{
		ageKeyFile: ageKeyFile,
	}
}

func (d *SOPSDecrypter) Decrypt(path string) ([]byte, error) {
	cmd := exec.Command("sops", "-d", path)

	cmd.Env = append(
		cmd.Environ(),
		"SOPS_AGE_KEY_FILE="+d.ageKeyFile,
	)

	return cmd.Output()
}
