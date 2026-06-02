package renderer

type Renderer interface {
    Render(spec model.ContainerSpec) (Unit, error)
}
