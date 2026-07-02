package normalize

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/juli3nk/podcd/internal/model"
)

type normalizedVolume struct {
	Name string

	Driver  string
	Options []string

	Labels []string
}

func defaultVolumeDriver(driver string) string {
	if driver == "" {
		return "local"
	}
	return driver
}

func normalizeVolume(v model.Volume) normalizedVolume {
	return normalizedVolume{
		Name:    v.Name,
		Driver:  defaultVolumeDriver(v.Driver),
		Options: normalizeMap(v.Options),
		Labels:  normalizeMap(v.Labels),
	}
}

func HashVolume(m model.Volume) string {
	n := normalizeVolume(m)

	data, _ := json.Marshal(n)

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
