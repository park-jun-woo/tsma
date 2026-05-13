//ff:func feature=graph type=helper control=selection
//ff:what Returns SkipDir for directories excluded from TS call graph analysis
package graph

import "path/filepath"

// skipTSBuildDir returns SkipDir for directories excluded from TS analysis.
func skipTSBuildDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "node_modules", "dist", "build", ".git", ".tsma":
		return filepath.SkipDir
	}
	return nil
}
