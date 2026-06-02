package runtime

type Runtime interface {
    Detect() RuntimeType
}
