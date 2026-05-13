//ff:func feature=index type=helper control=selection
//ff:what Returns SkipDir for directories excluded from Go source indexing
package index

import "path/filepath"

// skipGoDir returns SkipDir for directories that should be excluded from Go indexing.
func skipGoDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "vendor", ".git", ".tsma", "node_modules":
		return filepath.SkipDir
	}
	return nil
}
