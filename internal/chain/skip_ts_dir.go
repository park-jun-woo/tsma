//ff:func feature=chain type=helper control=selection
//ff:what Returns SkipDir for directories that should be excluded from TS search
package chain

import "path/filepath"

// skipTSDir returns SkipDir for directories that should be excluded from TS search.
func skipTSDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "node_modules", ".git", ".tsma", "dist", "build":
		return filepath.SkipDir
	}
	return nil
}
