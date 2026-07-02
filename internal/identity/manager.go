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

func New(basePath string) *IdentityManager {
	return &IdentityManager{
		basePath: basePath,
	}
}
