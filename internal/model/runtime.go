package model

type RuntimeType string

const (
    RuntimeDocker RuntimeType = "docker"
    RuntimePodman RuntimeType = "podman"
)

func (r RuntimeType) IsValid() bool {
    switch r {
    case RuntimeDocker, RuntimePodman:
        return true
    default:
        return false
    }
}
