//ff:func feature=chain type=helper control=sequence
//ff:what Checks if a file has a TS/JS extension suitable for source analysis
package chain

import "strings"

// isTSOrJSSourceFile checks if a file has a TS/JS extension suitable for source analysis.
func isTSOrJSSourceFile(path string) bool {
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
