//ff:func feature=index type=helper control=selection
//ff:what Returns SkipDir for directories excluded from TS/JS source indexing
package index

import "path/filepath"

// skipTSDir returns SkipDir for directories that should be excluded from TS indexing.
func skipTSDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "node_modules", "dist", "build", ".git", ".tsma":
		return filepath.SkipDir
	}
	return nil
}
