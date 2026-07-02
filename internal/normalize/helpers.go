package normalize

import (
	"fmt"
	"sort"
)

func normalizeMap(m map[string]string) []string {
	if m == nil {
		return []string{}
	}

	var out []string
	for k, v := range m {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(out)
	return out
}
