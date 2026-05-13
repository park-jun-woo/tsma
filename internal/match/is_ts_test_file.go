//ff:func feature=match type=helper control=iteration dimension=1
//ff:what Checks if a filename has a test/spec suffix
package match

import "strings"

// isTSTestFile returns true if the filename matches a TS/JS test file pattern.
func isTSTestFile(name string) bool {
	for _, suffix := range tsTestSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}
