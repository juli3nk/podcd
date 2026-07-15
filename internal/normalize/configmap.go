package normalize

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/juli3nk/podcd/internal/model"
)

type normalizedConfigMap struct {
	Name string

	Data string

	Labels []string
}

func normalizeConfigMap(v model.ConfigMap) normalizedConfigMap {
	return normalizedConfigMap{
		Name:   v.Name,
		Data:   v.Data,
		Labels: normalizeMap(v.Labels),
	}
}

func HashConfigMap(m model.ConfigMap) string {
	n := normalizeConfigMap(m)

	data, _ := json.Marshal(n)

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
