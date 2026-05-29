//ff:func feature=cli type=helper control=sequence
//ff:what Returns the mtime of a single test file as a time.Time, zero if missing
package cli

import (
	"os"
	"path/filepath"
	"time"
)

// testMtimeOf returns the modification time of a single test file (resolved
// relative to root). A missing or unstattable file yields the zero time.
func testMtimeOf(root, testFile string) time.Time {
	info, err := os.Stat(filepath.Join(root, testFile))
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
