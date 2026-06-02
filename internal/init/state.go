package init

type InitTask struct {
    Name    string
    Command []string
    Image   string
    Once    bool
}
