package wildcardstring

import (
	"regexp"
	"strings"
)

type WildcardString struct {
	Pattern string
}

func NewWildcardString(value string) *WildcardString {
	return &WildcardString{
		Pattern: value,
	}
}

func (w *WildcardString) Matches(value string) bool {
	result, _ := regexp.MatchString(wildCardToRegexp(w.Pattern), value)

	return result
}

// MatchesPath adds a wildcard to the end of the string to match any sub-paths.
func (w *WildcardString) MatchesPath(value string) bool {
	var pattern = w.Pattern

	pattern = wildCardToRegexp(pattern)

	if strings.HasSuffix(pattern, "/$") { // folders
		pattern = strings.TrimSuffix(pattern, "/$")
		pattern += "(/.*)?$"
	} else if strings.HasSuffix(pattern, "/.*$") {
		pattern = strings.TrimSuffix(pattern, "/.*$")

		pattern += "(/.*)?$"
	}

	result, _ := regexp.MatchString(pattern, value)

	return result
}

// wildCardToRegexp converts a wildcard Pattern to a regular expression Pattern.
func wildCardToRegexp(pattern string) string {
	components := strings.Split(pattern, "*")
	if len(components) == 1 {
		// if len is 1, there are no *'s, return exact match pattern
		return "^" + regexp.QuoteMeta(pattern) + "$"
	}

	var result strings.Builder

	for i, literal := range components {
		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}

		// Quote any regular expression meta characters in the literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}

	return "^" + result.String() + "$"
}
