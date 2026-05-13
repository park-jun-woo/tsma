//ff:func feature=index type=helper control=sequence
//ff:what Returns true if the path is a non-test Python source file eligible for indexing
package index

import "strings"

// isPySource returns true if the path is a Python source file eligible for indexing.
func isPySource(path string) bool {
	if !strings.HasSuffix(path, ".py") {
		return false
	}
	base := path
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		base = path[idx+1:]
	}
	if strings.HasPrefix(base, "test_") {
		return false
	}
	return true
}
