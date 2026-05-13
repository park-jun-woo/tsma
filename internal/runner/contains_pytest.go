//ff:func feature=runner type=helper control=sequence
//ff:what Checks if a file contains a pytest-related string
package runner

import (
	"os"
	"strings"
)

// containsPytest checks if a file contains a pytest-related string.
func containsPytest(path, pattern string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	content := strings.ToLower(string(data))
	return strings.Contains(content, strings.ToLower(pattern)) || strings.Contains(content, "pytest")
}
