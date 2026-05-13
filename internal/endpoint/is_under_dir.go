//ff:func feature=endpoint type=helper control=sequence
//ff:what Checks if a relative path is under the given directory
package endpoint

import (
	"path/filepath"
	"strings"
)

// isUnderDir checks if relPath is under the given directory.
func isUnderDir(relPath, dir string) bool {
	normalized := filepath.ToSlash(relPath)
	prefix := filepath.ToSlash(dir) + "/"
	srcPrefix := "src/" + prefix
	return strings.HasPrefix(normalized, prefix) || strings.HasPrefix(normalized, srcPrefix)
}
