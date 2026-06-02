package systemd

type Manager interface {
    Apply(unit Unit) error
    Reload() error
    Restart(name string) error
    Enable(name string) error
    Disable(name string) error
}
