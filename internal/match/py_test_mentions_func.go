//ff:func feature=match type=helper control=sequence
//ff:what Scans Python test file content for def test_funcname patterns
package match

import (
	"os"
	"regexp"
	"strings"
)

// pyTestMentionsFunc reads a Python test file and checks for
// `def test_funcname` patterns.
func pyTestMentionsFunc(testFilePath string, funcName string) bool {
	data, err := os.ReadFile(testFilePath)
	if err != nil {
		return false
	}
	content := string(data)

	// Match def test_funcname (case-insensitive comparison).
	pattern := regexp.MustCompile(`def\s+test_` + regexp.QuoteMeta(strings.ToLower(funcName)) + `\b`)
	if pattern.MatchString(strings.ToLower(content)) {
		return true
	}

	// Also try exact match without lowercasing.
	exactPattern := regexp.MustCompile(`def\s+test_` + regexp.QuoteMeta(funcName) + `\b`)
	return exactPattern.MatchString(content)
}
