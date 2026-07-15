package normalize

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/juli3nk/podcd/internal/model"
)

type normalizedInitTask struct {
	Name string

	Image   string
	Command []string

	Volumes []model.VolumeRef
	Devices []string

	Networks []model.NetworkRef
	DNS      []string

	Env        []string
	ConfigMaps []model.FileRef
	Secrets    []model.FileRef

	Once bool
}

func normalizeInitTask(spec model.InitTask) normalizedInitTask {
	return normalizedInitTask{
		Name: spec.Name,

		Image:   spec.Image,
		Command: spec.Command,

		Volumes: sortVolumes(spec.Volumes),
		Devices: sortStrings(spec.Devices),

		Networks: sortNetworks(spec.Networks),
		DNS:      sortStrings(spec.DNS),

		Env:        normalizeEnv(spec.Env),
		ConfigMaps: normalizeFiles(spec.ConfigMaps),
		Secrets:    normalizeFiles(spec.Secrets),

		Once: spec.Once,
	}
}

func HashInitTask(m model.InitTask) string {
	n := normalizeInitTask(m)

	data, _ := json.Marshal(n)

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
