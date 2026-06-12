package renderer

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/juli3nk/podcd/internal/model"
)

type DockerRenderer struct{}

var dockerTemplate = `
[Unit]
Description={{.Name}} (Docker)
After=network.target

[Service]
Restart=always
ExecStart=/usr/bin/{{.CommandLine}}
ExecStop=/usr/bin/docker stop {{ .Name }}
ExecStopPost=/usr/bin/docker rm {{ .Name }}

[Install]
WantedBy=multi-user.target
`

func (r *DockerRenderer) Render(spec model.ContainerSpec) ([]Unit, error) {
	data := struct {
		Name        string
		CommandLine string
	}{
		Name:        spec.Name,
		CommandLine: buildDockerCommand(spec),
	}
	tmpl, err := template.New("docker").Parse(dockerTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	unit := Unit{
		Name:    spec.Name + ".service",
		Path:    "/etc/systemd/system/" + spec.Name + ".service",
		Content: buf.String(),
	}

	return []Unit{unit}, nil
}

func buildDockerCommand(spec model.ContainerSpec) string {
	args := []string{"docker", "run", "--rm", "--name", spec.Name}

	for _, p := range spec.Ports {
		args = append(args, fmt.Sprintf("-p %d:%d", p.HostPort, p.ContainerPort))
	}

	for k, v := range spec.Env {
		args = append(args, fmt.Sprintf("-e %s=%s", k, v))
	}

	args = append(args, spec.Image)
	args = append(args, spec.Command...)

	return strings.Join(args, " ")
}
