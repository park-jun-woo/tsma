//ff:func feature=coverage type=helper control=sequence
//ff:what Checks if a coverage key matches the target relative path by direct or suffix comparison
package coverage

import (
	"path/filepath"
	"strings"
)

// matchesCoverageKey checks if a coverage key matches the target relative path.
func matchesCoverageKey(key, normalizedRel, relFile, projectRoot string) bool {
	normalizedKey := filepath.ToSlash(key)

	if normalizedKey == normalizedRel {
		return true
	}

	if strings.HasSuffix(normalizedKey, "/"+normalizedRel) || strings.HasSuffix(normalizedKey, normalizedRel) {
		return true
	}

	if projectRoot == "" {
		return false
	}

	absFile := filepath.ToSlash(filepath.Join(projectRoot, relFile))
	return normalizedKey == absFile
}
