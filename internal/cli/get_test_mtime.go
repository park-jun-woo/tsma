//ff:func feature=cli type=helper control=sequence
//ff:what Returns the mtime of a test file in RFC3339 format
package cli

import (
	"os"
	"path/filepath"
	"time"
)

// getTestMtime returns the mtime of a test file in RFC3339 format.
func getTestMtime(root, testFile string) string {
	abs := filepath.Join(root, testFile)
	info, err := os.Stat(abs)
	if err != nil {
		return ""
	}
	return info.ModTime().Format(time.RFC3339)
}
