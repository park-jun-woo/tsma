//ff:func feature=coverage type=helper control=sequence
//ff:what Checks if a Python coverage path matches a target relative path
package coverage

import (
	"path/filepath"
	"strings"
)

// matchesPyPath checks if a coverage path matches a target relative path.
func matchesPyPath(covPath, targetFile, projectRoot string) bool {
	if covPath == targetFile {
		return true
	}

	absTarget := filepath.Join(projectRoot, targetFile)
	if covPath == absTarget {
		return true
	}

	normalized := filepath.ToSlash(covPath)
	normalizedTarget := filepath.ToSlash(targetFile)
	if strings.HasSuffix(normalized, normalizedTarget) {
		idx := len(normalized) - len(normalizedTarget)
		if idx == 0 || normalized[idx-1] == '/' {
			return true
		}
	}
	return false
}
