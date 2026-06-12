package renderer

import (
	"fmt"

	"github.com/juli3nk/podcd/internal/model"
)

type Unit struct {
	Name    string // ex: nginx.service ou nginx.container
	Content string
	Path    string
}

type Renderer interface {
	Render(spec model.ContainerSpec) ([]Unit, error)
}

func New(rt model.RuntimeType) (Renderer, error) {
	switch rt {
	case model.RuntimeDocker:
		return &DockerRenderer{}, nil
	case model.RuntimePodman:
		return &QuadletRenderer{}, nil
	default:
		return nil, fmt.Errorf("unsupported runtime: %s", rt)
	}
}
