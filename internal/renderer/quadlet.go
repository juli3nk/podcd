package renderer

import (
	"bytes"
	"text/template"

	"github.com/juli3nk/podcd/internal/model"
)

type QuadletRenderer struct{}

var quadletTemplate = `
[Unit]
Description={{ .Name }} (Podman)

[Container]
Image={{ .Image }}
{{- range .Ports }}
PublishPort={{.}}
{{- end }}
{{- range $k, $v := .Env}}
Environment={{$k}}={{$v}}
{{- end }}

[Service]
Restart=always

[Install]
WantedBy=multi-user.target
`

func (r *QuadletRenderer) Render(spec model.ContainerSpec) ([]Unit, error) {
	tmpl, err := template.New("podman").Parse(quadletTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, spec); err != nil {
		return nil, err
	}

	unit := Unit{
		Name:    spec.Name + ".service",
		Path:    "/etc/containers/systemd/" + spec.Name + ".container",
		Content: buf.String(),
	}

	return []Unit{unit}, nil
}
