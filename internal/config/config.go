package config

import (
	"github.com/juli3nk/go-utils/user"
)

type Config struct {
	Source struct {
		RepoURL string `yaml:"repoUrl"`
		Version string `yaml:"version"`
	} `yaml:"source"`

	Runtime struct {
		Type string `yaml:"type"` // auto | docker | podman
	} `yaml:"runtime"`

	Controller struct {
		Interval string `yaml:"interval"` // 10s
	} `yaml:"controller"`

	Server struct {
		LogLevel string `yaml:"logLevel"`
	} `yaml:"server"`
}

func IsUserMode() bool {
	u := user.New()

	if u.IsRoot() {
		return false
	}

	return true
}
