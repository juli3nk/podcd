package normalize

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/juli3nk/podcd/internal/model"
)

type normalizedContainer struct {
	Name string

	Image   string
	Command []string

	Volumes []model.VolumeRef
	Devices []string

	Networks []model.NetworkRef
	DNS      []string
	Ports    []model.PortSpec

	Env        []string
	ConfigMaps []model.FileRef
	Secrets    []model.FileRef

	AddCapabilities  []string
	DropCapabilities []string
}

func normalizeEnv(env map[string]string) []string {
	var out []string
	for k, v := range env {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(out)
	return out
}

func sortStrings(in []string) []string {
	out := append([]string{}, in...)
	sort.Strings(out)
	return out
}

func sortPorts(p []model.PortSpec) []model.PortSpec {
	out := append([]model.PortSpec{}, p...)

	sort.Slice(out, func(i, j int) bool {
		if out[i].ContainerPort == out[j].ContainerPort {
			return out[i].HostPort < out[j].HostPort
		}
		return out[i].ContainerPort < out[j].ContainerPort
	})

	return out
}

func sortVolumes(v []model.VolumeRef) []model.VolumeRef {
	out := append([]model.VolumeRef{}, v...)

	sort.Slice(out, func(i, j int) bool {
		if out[i].Source == out[j].Source {
			return out[i].Target < out[j].Target
		}
		return out[i].Source < out[j].Source
	})

	return out
}

func sortNetworks(n []model.NetworkRef) []model.NetworkRef {
	out := append([]model.NetworkRef{}, n...)

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})

	return out
}

func normalizeFiles(files []model.FileRef) []model.FileRef {
	result := slices.Clone(files)

	slices.SortFunc(result, func(a, b model.FileRef) int {
		if a.Name != b.Name {
			return strings.Compare(a.Name, b.Name)
		}

		return strings.Compare(a.Target, b.Target)
	})

	return result
}

func normalizeContainer(spec model.Container) normalizedContainer {
	return normalizedContainer{
		Name: spec.Name,

		Image:   spec.Image,
		Command: spec.Command,

		Volumes: sortVolumes(spec.Volumes),
		Devices: sortStrings(spec.Devices),

		Networks: sortNetworks(spec.Networks),
		DNS:      sortStrings(spec.DNS),
		Ports:    sortPorts(spec.Ports),

		Env:        normalizeEnv(spec.Env),
		ConfigMaps: normalizeFiles(spec.ConfigMaps),
		Secrets:    normalizeFiles(spec.Secrets),

		AddCapabilities:  sortStrings(spec.AddCapabilities),
		DropCapabilities: sortStrings(spec.DropCapabilities),
	}
}

func HashContainer(m model.Container) string {
	n := normalizeContainer(m)

	data, _ := json.Marshal(n)

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
