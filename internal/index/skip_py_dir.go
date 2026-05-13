//ff:func feature=index type=helper control=selection
//ff:what Returns SkipDir for directories excluded from Python source indexing
package index

import "path/filepath"

// skipPyDir returns SkipDir for directories that should be excluded from Python indexing.
func skipPyDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "__pycache__", ".venv", "venv", ".git", ".tsma":
		return filepath.SkipDir
	}
	return nil
}
