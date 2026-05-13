//ff:func feature=graph type=helper control=sequence
//ff:what Extracts the package directory from a file path
package graph

import "path/filepath"

// pkgDirOf returns the directory portion of a file path, or empty string for root.
func pkgDirOf(file string) string {
	dir := filepath.Dir(file)
	if dir == "." {
		return ""
	}
	return dir
}
