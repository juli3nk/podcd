package model

type ContainerSpec struct {
    Name    string
    Image   string
    Ports   []string
    Env     map[string]string
    Runtime string
}
