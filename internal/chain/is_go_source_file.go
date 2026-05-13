//ff:func feature=chain type=helper control=sequence
//ff:what Returns true if the path is a non-test Go source file
package chain

import "strings"

// isGoSourceFile returns true if the path is a non-test Go source file.
func isGoSourceFile(path string) bool {
	return strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go")
}
