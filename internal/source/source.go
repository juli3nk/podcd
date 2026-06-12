package source

type ChangeType string

const (
	Added    ChangeType = "added"
	Modified ChangeType = "modified"
	Deleted  ChangeType = "deleted"
)

type Change struct {
	Path string
	Type ChangeType
}

type Source interface {
	Fetch() error

	HasChanged() (bool, error) // simple trigger
	Diff() ([]Change, error)   // granularité fin

	Revision() string
	Path() string
}
