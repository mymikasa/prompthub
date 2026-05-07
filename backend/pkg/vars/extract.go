package vars

import (
	"regexp"
	"sort"
	"strings"
)

var varPattern = regexp.MustCompile(`\{\{([a-zA-Z_][a-zA-Z0-9_]*)\}\}`)

func Extract(text string) []string {
	matches := varPattern.FindAllStringSubmatch(text, -1)
	seen := make(map[string]struct{})
	for _, m := range matches {
		seen[m[1]] = struct{}{}
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func IsValidName(name string) bool {
	if len(name) == 0 || len(name) > 100 {
		return false
	}
	return varPattern.MatchString("{{" + name + "}}")
}

func NormalizeName(name string) string {
	return strings.TrimSpace(name)
}
