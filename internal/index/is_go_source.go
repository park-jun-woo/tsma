//ff:func feature=index type=helper control=sequence
//ff:what Returns true if the path is a non-test non-mock non-generated Go source file
package index

import "strings"

// isGoSource returns true if the path is a Go source file eligible for indexing.
func isGoSource(path string) bool {
	if !strings.HasSuffix(path, ".go") {
		return false
	}
	if strings.HasSuffix(path, "_test.go") {
		return false
	}
	if strings.HasSuffix(path, "_gen.go") || strings.HasSuffix(path, ".gen.go") || strings.HasSuffix(path, ".pb.go") {
		return false
	}
	base := path
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		base = path[idx+1:]
	}
	if strings.HasPrefix(base, "mock_") {
		return false
	}
	return true
}
