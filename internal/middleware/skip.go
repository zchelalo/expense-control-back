package middleware

import "strings"

var excluded = map[string]struct{}{
  "/api/metrics": {},
  "/api/health":  {},
}

func ShouldInstrument(path string) bool {
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}
	_, excluded := excluded[path]
	return !excluded
}