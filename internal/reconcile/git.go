package reconcile

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/juli3nk/podcd/internal/model"

	yaml "go.yaml.in/yaml/v4"
)

func FindYAMLFiles(basePath string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case ".yaml", ".yml":
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func loadFromGit(basePath string) (model.RootSpec, error) {
	var result model.RootSpec

	files, err := FindYAMLFiles(basePath)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		fmt.Println(file)
		data, err := os.ReadFile(file)
		if err != nil {
			return result, err
		}

		var meta struct {
			Kind string `yaml:"kind"`
		}

		if err := yaml.Unmarshal(data, &meta); err != nil {
			return result, err
		}

		switch strings.ToLower(meta.Kind) {
		case string(ResourceNetwork):
			var n model.Network
			yaml.Unmarshal(data, &n)
			result.Networks = append(result.Networks, n)
		case string(ResourceVolume):
			var v model.Volume
			yaml.Unmarshal(data, &v)
			result.Volumes = append(result.Volumes, v)
		case string(ResourceSecret):
			var s model.Secret
			yaml.Unmarshal(data, &s)
			s.Filepath = file
			result.Secrets = append(result.Secrets, s)
		case string(ResourceContainer):
			var c model.Container
			yaml.Unmarshal(data, &c)
			result.Containers = append(result.Containers, c)
		}
	}

	return result, nil
}
