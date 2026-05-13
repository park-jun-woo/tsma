//ff:func feature=index type=helper control=sequence
//ff:what Returns true if the path is a non-test TS/JS source file eligible for indexing
package index

import "strings"

// isTSSource returns true if the path is a TS/JS source file eligible for indexing.
func isTSSource(path string) bool {
	if strings.HasSuffix(path, ".d.ts") {
		return false
	}
	if strings.HasSuffix(path, ".test.ts") || strings.HasSuffix(path, ".test.js") {
		return false
	}
	if strings.HasSuffix(path, ".spec.ts") || strings.HasSuffix(path, ".spec.js") {
		return false
	}
	return strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".js")
}
