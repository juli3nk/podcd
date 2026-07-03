package normalize

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/juli3nk/podcd/internal/model"
)

type normalizedSecret struct {
	Name string

	EncryptedData []string

	Labels []string
}

func normalizeSecret(v model.Secret) normalizedSecret {
	return normalizedSecret{
		Name:   v.Name,
		Labels: normalizeMap(v.Labels),
	}
}

func HashSecret(m model.Secret) string {
	n := normalizeSecret(m)

	data, _ := json.Marshal(n)

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
