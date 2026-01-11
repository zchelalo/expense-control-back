package middleware

import "strings"

var excluded = map[string]struct{}{
	"/metrics": {},
	"/health":  {},
}

func ShouldSkipObservability(path string) bool {
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}
	_, skip := excluded[path]
	return !skip
}