package renderer

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/juli3nk/podcd/internal/model"
	"github.com/juli3nk/podcd/internal/runtime"
)

type SystemdRenderer struct {
	backend    runtime.Backend
	binaryPath string
}

var systemdUnitTemplate = `
[Unit]
Description={{ .Name }} ({{ .Backend }})
After=network.target

[Service]
Restart=always
ExecStart={{ .BinaryPath }}/{{ .CommandLine }}
ExecStop={{ .BinaryPath }} container stop {{ .Name }}
ExecStopPost={{ .BinaryPath }} container rm {{ .Name }}

[Install]
WantedBy=multi-user.target
`

func (r *SystemdRenderer) Render(spec model.Container, hash string) ([]Unit, error) {
	data := struct {
		Backend     string
		BinaryPath  string
		Name        string
		CommandLine string
	}{
		Backend:     string(r.backend),
		BinaryPath:  r.binaryPath,
		Name:        spec.Name,
		CommandLine: r.buildRunCommand(spec, hash),
	}
	tmpl, err := template.New(data.Backend).Parse(systemdUnitTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	unit := Unit{
		Name:    spec.Name + ".service",
		Content: buf.String(),
	}

	return []Unit{unit}, nil
}

func (r *SystemdRenderer) buildRunCommand(spec model.Container, hash string) string {
	switch r.backend {
	case runtime.BackendDocker:
		args := runtime.BuildDockerContainerRunArgs(spec, hash)
		return strings.Join(args, " ")

	case runtime.BackendPodman:
		args := runtime.BuildDockerContainerRunArgs(spec, hash)
		return strings.Join(args, " ")

	default:
		panic("unsupported backend")
	}
}
