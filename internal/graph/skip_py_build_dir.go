//ff:func feature=graph type=helper control=selection
//ff:what Returns SkipDir for directories excluded from Python call graph analysis
package graph

import "path/filepath"

// skipPyBuildDir returns SkipDir for directories excluded from Python analysis.
func skipPyBuildDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "__pycache__", ".venv", "venv", ".git", ".tsma":
		return filepath.SkipDir
	}
	return nil
}
