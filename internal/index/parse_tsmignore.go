//ff:func feature=index type=loader control=iteration dimension=1
//ff:what .tsmignore 파일을 읽어 패턴 목록을 반환 (없으면 nil)
package index

import (
	"bufio"
	"os"
	"strings"
)

// ParseTsmIgnore reads a .tsmignore file and returns the list of patterns.
// Returns nil if the file doesn't exist.
func ParseTsmIgnore(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var patterns []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns
}
