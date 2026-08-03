package config

import (
	"fmt"
	"net/url"
	"time"

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

func (c Config) Validate() error {
	if c.Source.RepoURL == "" {
		return fmt.Errorf("source.repoUrl is required")
	}
	if _, err := url.Parse(c.Source.RepoURL); err != nil {
		return fmt.Errorf("source.repoUrl is invalid: %w", err)
	}
	if c.Controller.Interval != "" {
		if _, err := time.ParseDuration(c.Controller.Interval); err != nil {
			return fmt.Errorf("controller.interval is invalid: %w", err)
		}
	}
	if c.Runtime.Type != "" && c.Runtime.Type != "auto" && c.Runtime.Type != "docker" && c.Runtime.Type != "podman" {
		return fmt.Errorf("runtime.type must be auto, docker, or podman")
	}
	return nil
}
