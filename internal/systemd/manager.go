package systemd

import (
	"fmt"
	"os/exec"
)

type Manager interface {
	Reload() error

	Enable(name string) error
	EnableNow(name string) error
	Disable(name string) error
	DisableNow(name string) error

	Start(name string) error
	Stop(name string) error
	Restart(name string) error
}

type Systemd struct {
	UserMode bool
}

func New(userMode bool) Manager {
	return &Systemd{UserMode: userMode}
}

func (s *Systemd) Reload() error {
	return s.run("daemon-reload")
}

func (s *Systemd) Enable(name string) error {
	return s.run("enable", name)
}

func (s *Systemd) EnableNow(name string) error {
	return s.run("enable", "--now", name)
}

func (s *Systemd) Disable(name string) error {
	return s.run("disable", name)
}

func (s *Systemd) DisableNow(name string) error {
	return s.run("disable", "--now", name)
}

func (s *Systemd) Start(name string) error {
	return s.run("start", name)
}

func (s *Systemd) Stop(name string) error {
	return s.run("stop", name)
}

func (s *Systemd) Restart(name string) error {
	return s.run("restart", name)
}

func (s *Systemd) run(args ...string) error {
	finalArgs := []string{}

	if s.UserMode {
		finalArgs = append(finalArgs, "--user")
	}

	finalArgs = append(finalArgs, args...)

	cmd := exec.Command("systemctl", finalArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("systemctl %v failed: %v\noutput: %s",
			finalArgs, err, string(output))
	}

	return nil
}
