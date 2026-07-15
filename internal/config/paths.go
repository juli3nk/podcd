package config

import (
	"os"
	"path/filepath"
)

type Paths struct {
	Config    string
	State     string
	Workspace string

	IdentityDir       string
	RuntimeStorageDir string

	Socket string
}

func DefaultPaths(userMode bool) Paths {
	appName := "podcd"

	if userMode {
		home, _ := os.UserHomeDir()

		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			configDir = filepath.Join(home, ".config")
		}

		dataDir := os.Getenv("XDG_DATA_HOME")
		if dataDir == "" {
			dataDir = filepath.Join(home, ".local", "share")
		}

		runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
		if runtimeDir == "" {
			uid := ""
			runtimeDir = filepath.Join("/run", "user", uid)
		}

		return Paths{
			Config:            filepath.Join(configDir, appName, "config.yaml"),
			State:             filepath.Join(dataDir, appName, "state.json"),
			Workspace:         filepath.Join(dataDir, appName, "gitops-repo"),
			IdentityDir:       filepath.Join(configDir, appName, "identity"),
			RuntimeStorageDir: filepath.Join(dataDir, appName, "runtime"),
			Socket:            filepath.Join(runtimeDir, appName, "podcd.sock"),
		}
	}

	return Paths{
		Config:            filepath.Join("/etc", appName, "config.yaml"),
		State:             filepath.Join("/var/lib", appName, "state.json"),
		Workspace:         filepath.Join("/var/lib", appName, "gitops-repo"),
		IdentityDir:       filepath.Join("/etc", appName, "identity"),
		RuntimeStorageDir: filepath.Join("/run", appName, "runtime"),
		Socket:            filepath.Join("/run", appName, "podcd.sock"),
	}
}

func EnsureDirectories(paths Paths) error {
	dirs := []string{
		filepath.Dir(paths.Config),
		filepath.Dir(paths.Socket),
		filepath.Dir(paths.State),
		paths.Workspace,
		paths.IdentityDir,
		paths.RuntimeStorageDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
