package reconciler

type Reconciler struct {
    Runtime   runtime.Runtime
    Renderer  renderer.Renderer
    Systemd   systemd.Manager
}
