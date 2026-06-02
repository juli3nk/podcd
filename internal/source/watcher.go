package source

type Source interface {
    Fetch() error
    Diff() ([]Change, error)
}
