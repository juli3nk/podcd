package state

type Store interface {
    IsDone(name string) bool
    MarkDone(name string) error
    Save() error
    Load() error
}
