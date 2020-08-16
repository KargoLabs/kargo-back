package normalize

import "strings"

// NormalizeName makes all letter uppercase and removes additional spaces
func NormalizeName(name string) string {
	return strings.ToUpper(NormalizeSpaces(name))
}

// NormalizeSpaces removes additionals white spaces
func NormalizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
