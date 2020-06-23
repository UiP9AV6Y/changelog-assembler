package change

import (
	"strings"
)

const (
	AnnotationSeparator = '='
)

// ParseAnnotations creates a map from key/value pairs.
func ParseAnnotations(list []string) map[string]string {
	var key, value string
	result := make(map[string]string, len(list))

	for _, item := range list {
		i := strings.IndexByte(item, AnnotationSeparator)

		if i < 0 {
			key = item
			value = ""
		} else {
			key = item[:i]
			value = item[i+1:]
		}

		result[key] = value
	}

	return result
}
