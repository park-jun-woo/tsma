//ff:func feature=graph type=helper control=selection
//ff:what Returns SkipDir for directories excluded from Go call graph analysis
package graph

import "path/filepath"

// skipGoBuildDir returns SkipDir for directories excluded from Go analysis.
func skipGoBuildDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "vendor", ".git", ".tsma", "node_modules":
		return filepath.SkipDir
	}
	return nil
}
