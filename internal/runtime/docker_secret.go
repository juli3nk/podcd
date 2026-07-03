package runtime

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/juli3nk/podcd/internal/model"
)

func (r *DockerRuntime) ListSecrets(filter Labels) ([]SecretInfo, error) {
	var result []SecretInfo

	entries, err := os.ReadDir(r.basePath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".meta.json" {
			continue
		}

		labels, err := r.loadSecretLabels(
			filepath.Join(r.basePath, entry.Name()),
		)
		if err != nil {
			return nil, err
		}

		result = append(result, SecretInfo{Name: entry.Name(), Labels: labels})
	}

	return result, nil
}

func (r *DockerRuntime) CreateSecret(secret model.Secret, data []byte, hash string) error {
	if err := os.MkdirAll(r.basePath, 0700); err != nil {
		return err
	}

	path := filepath.Join(r.basePath, secret.Name)

	return os.WriteFile(
		path,
		data,
		0600,
	)
}

func (r *DockerRuntime) RemoveSecret(name string) error {
	path := filepath.Join(r.basePath, name)

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (r *DockerRuntime) loadSecretLabels(secretFilename string) (Labels, error) {
	metaFile := secretFilename + ".meta.json"

	data, err := os.ReadFile(metaFile)
	if err != nil {
		return nil, err
	}

	var labels Labels

	if err := json.Unmarshal(data, &labels); err != nil {
		return nil, err
	}

	return labels, nil
}
