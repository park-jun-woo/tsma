//ff:func feature=runner type=implementation control=iteration dimension=1
//ff:what Parses a Go test file and returns test function names
package runner

import "strings"

// ExtractTestFuncs parses a Go test file and returns test function names.
func ExtractTestFuncs(filePath string) ([]string, error) {
	data, err := readFileBytes(filePath)
	if err != nil {
		return nil, err
	}
	content := string(data)
	var funcs []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "func Test") {
			continue
		}
		rest := line[len("func "):]
		idx := strings.IndexByte(rest, '(')
		if idx > 0 {
			name := rest[:idx]
			funcs = append(funcs, name)
		}
	}
	return funcs, nil
}
