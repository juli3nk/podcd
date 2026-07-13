package identity

type Manager interface {
	EnsureSSHKey() error
	EnsureAgeKey() error

	SSHPublicKey() (string, error)
	AgePublicKey() (string, error)
}

type IdentityManager struct {
	basePath string
}

func New(basePath string) Manager {
	return &IdentityManager{
		basePath: basePath,
	}
}
