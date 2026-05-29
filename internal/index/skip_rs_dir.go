//ff:func feature=index type=helper control=selection
//ff:what Returns SkipDir for directories excluded from Rust source indexing
package index

import "path/filepath"

// skipRsDir returns SkipDir for directories that should be excluded from Rust indexing.
func skipRsDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "target", ".git", ".tsma":
		return filepath.SkipDir
	}
	return nil
}
