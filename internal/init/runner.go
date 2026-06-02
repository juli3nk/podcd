package init

type InitRunner interface {
    Run(task InitTask) error
    AlreadyDone(task InitTask) bool
}
