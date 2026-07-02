package normalize

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/juli3nk/podcd/internal/model"
)

type normalizedIPAM struct {
	Driver  string
	Options []string
	Config  []model.IPAMConfig
}

type normalizedNetwork struct {
	Name string

	Driver  string
	Options []string

	IPAM normalizedIPAM

	Internal   bool
	Attachable bool

	Labels []string
}

func defaultNetworkDriver(driver string) string {
	if driver == "" {
		return "bridge"
	}
	return driver
}

func defaultIPAMDriver(driver string) string {
	if driver == "" {
		return "default"
	}
	return driver
}

func normalizeIPAMConfig(in []model.IPAMConfig) []model.IPAMConfig {
	out := append([]model.IPAMConfig{}, in...)

	sort.Slice(out, func(i, j int) bool {
		return out[i].Subnet < out[j].Subnet
	})

	return out
}

func normalizeIPAM(ipam model.IPAM) normalizedIPAM {
	return normalizedIPAM{
		Driver:  defaultIPAMDriver(ipam.Driver),
		Options: normalizeMap(ipam.Options),
		Config:  normalizeIPAMConfig(ipam.Config),
	}
}

func normalizeNetwork(n model.Network) normalizedNetwork {
	return normalizedNetwork{
		Name:       n.Name,
		Driver:     defaultNetworkDriver(n.Driver),
		Options:    normalizeMap(n.Options),
		IPAM:       normalizeIPAM(n.IPAM),
		Internal:   n.Internal,
		Attachable: n.Attachable,
		Labels:     normalizeMap(n.Labels),
	}
}

func HashNetwork(m model.Network) string {
	n := normalizeNetwork(m)

	data, _ := json.Marshal(n)

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
