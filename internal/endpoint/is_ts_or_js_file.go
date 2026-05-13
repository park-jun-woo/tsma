//ff:func feature=endpoint type=helper control=sequence
//ff:what Checks if a file has a .ts or .js extension excluding .d.ts
package endpoint

import "strings"

// isTSOrJSFile checks if the file has a .ts or .js extension (excluding .d.ts).
func isTSOrJSFile(path string) bool {
	if strings.HasSuffix(path, ".d.ts") {
		return false
	}
	return strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".js")
}
