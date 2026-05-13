//ff:func feature=match type=helper control=sequence
//ff:what Scans test file content for describe/test/it calls mentioning the function name
package match

import (
	"os"
	"regexp"
)

// tsTestMentionsFunc reads a test file and checks for describe/test/it
// calls that mention the function name.
func tsTestMentionsFunc(testFilePath string, funcName string) bool {
	data, err := os.ReadFile(testFilePath)
	if err != nil {
		return false
	}
	content := string(data)

	pattern := regexp.MustCompile(`(?:describe|test|it)\s*\(\s*['"]` + regexp.QuoteMeta(funcName) + `['"]`)
	return pattern.MatchString(content)
}
