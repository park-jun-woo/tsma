//ff:func feature=detect type=helper control=sequence
//ff:what Checks if a file contains a pytest-related string (case-insensitive)
package detect

import (
	"os"
	"strings"
)

// containsPytest reports whether the file at path contains pattern
// (case-insensitive). A missing/unreadable file yields false. This mirrors the
// runner-package helper of the same name (cf. test_file_from_match.go's
// duplication precedent across packages); detect is the SSOT for pytest
// detection so the check body lives here.
func containsPytest(path, pattern string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	content := strings.ToLower(string(data))
	return strings.Contains(content, strings.ToLower(pattern))
}
