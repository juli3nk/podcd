package secrets

type SecretStore interface {
    Get(name string) (map[string]string, error)
}
