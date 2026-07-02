package reconcile

import (
	"sort"

	"github.com/juli3nk/podcd/internal/model"
)

func unionKeysContainer(
	actual map[string]RuntimeObject,
	desired map[string]model.Container,
) []string {
	set := make(map[string]struct{})

	for k := range actual {
		set[k] = struct{}{}
	}
	for k := range desired {
		set[k] = struct{}{}
	}

	var result []string
	for k := range set {
		result = append(result, k)
	}

	sort.Strings(result)
	return result
}

func unionKeysNetwork(
	actual map[string]RuntimeObject,
	desired map[string]model.Network,
) []string {
	set := make(map[string]struct{})

	for k := range actual {
		set[k] = struct{}{}
	}
	for k := range desired {
		set[k] = struct{}{}
	}

	var result []string
	for k := range set {
		result = append(result, k)
	}

	sort.Strings(result)
	return result
}

func unionKeysSecret(
	actual map[string]RuntimeObject,
	desired map[string]model.Secret,
) []string {
	set := make(map[string]struct{})

	for k := range actual {
		set[k] = struct{}{}
	}
	for k := range desired {
		set[k] = struct{}{}
	}

	var result []string
	for k := range set {
		result = append(result, k)
	}

	sort.Strings(result)
	return result
}

func unionKeysVolume(
	actual map[string]RuntimeObject,
	desired map[string]model.Volume,
) []string {
	set := make(map[string]struct{})

	for k := range actual {
		set[k] = struct{}{}
	}
	for k := range desired {
		set[k] = struct{}{}
	}

	var result []string
	for k := range set {
		result = append(result, k)
	}

	sort.Strings(result)
	return result
}
