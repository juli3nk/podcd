package runtime

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/juli3nk/podcd/internal/model"
)

func (r *PodmanRuntime) ListConfigMaps(filter Labels) ([]ConfigMapInfo, error) {
	var result []ConfigMapInfo

	entries, err := os.ReadDir(r.configMapDir())
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".meta.json" {
			continue
		}

		labels, err := r.loadConfigMapLabels(
			r.configMapPath(entry.Name()),
		)
		if err != nil {
			return nil, err
		}

		result = append(result, ConfigMapInfo{Name: entry.Name(), Labels: labels})
	}

	return result, nil
}

func (r *PodmanRuntime) CreateConfigMap(configMap model.ConfigMap, hash string) error {
	if err := os.MkdirAll(r.configMapDir(), 0700); err != nil {
		return err
	}

	path := r.configMapPath(configMap.Name)
	data := []byte(configMap.Data)

	return os.WriteFile(
		path,
		data,
		0600,
	)
}

func (r *PodmanRuntime) RemoveConfigMap(name string) error {
	path := r.configMapPath(name)

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (r *PodmanRuntime) loadConfigMapLabels(configMapFilename string) (Labels, error) {
	metaFile := configMapFilename + ".meta.json"

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
