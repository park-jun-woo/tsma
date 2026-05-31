//ff:func feature=cli type=helper control=sequence
//ff:what Returns the absolute package directory for a match (dir of its first test file)
package cli

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
)

// goPkgDirOf returns the absolute package directory for a match, derived from the
// directory of its first matched test file (all files in a match share a
// package). Relative match paths are resolved against the project root (not the
// process cwd) so batching works regardless of where tsma is invoked from.
func goPkgDirOf(root string, m match.TestMatch) string {
	f := m.Files[0]
	if !filepath.IsAbs(f) {
		f = filepath.Join(root, f)
	}
	return filepath.Dir(f)
}
