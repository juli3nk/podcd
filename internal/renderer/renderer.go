package renderer

import (
	"github.com/juli3nk/podcd/internal/model"
	"github.com/juli3nk/podcd/internal/runtime"
)

type Unit struct {
	Name    string // ex: nginx.service ou nginx.container
	Content string
}

type Renderer interface {
	Render(spec model.ContainerSpec, hash string) ([]Unit, error)
}

func New(rt runtime.Backend) (Renderer, error) {
	binaryPath, err := rt.BinaryPath()
	if err != nil {
		return nil, err
	}

	return &SystemdRenderer{
		backend:    rt,
		binaryPath: binaryPath,
	}, nil
}
