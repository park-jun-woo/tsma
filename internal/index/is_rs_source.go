//ff:func feature=index type=helper control=sequence
//ff:what Returns true if the path is a non-test Rust source file eligible for indexing
package index

import "strings"

// isRsSource returns true if the path is a Rust source file eligible for indexing.
func isRsSource(path string) bool {
	if !strings.HasSuffix(path, ".rs") {
		return false
	}
	base := path
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		base = path[idx+1:]
	}
	// build.rs is a build script, not library/binary source.
	if base == "build.rs" {
		return false
	}
	return true
}
