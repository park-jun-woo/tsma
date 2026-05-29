//ff:func feature=index type=implementation control=iteration dimension=1
//ff:what Parses a single TS/JS file and extracts function declarations line by line
package index

import (
	"bufio"
	"os"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// indexTSFile parses a single TS/JS file and extracts function declarations.
func indexTSFile(relPath, absPath string) []model.Function {
	f, err := os.Open(absPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	relDir := pkgDirOf(relPath)

	var functions []model.Function
	var currentClass string
	var classIndent int
	lineNum := 0
	lastNonEmptyLine := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed != "" {
			lastNonEmptyLine = lineNum
		}

		if m := tsClassPattern.FindStringSubmatch(trimmed); m != nil {
			closePrevTSEndLine(functions, lineNum, lastNonEmptyLine)
			currentClass = m[1]
			classIndent = countLeadingSpaces(line)
			continue
		}

		if currentClass != "" && resetTSClassContext(trimmed, line, classIndent) {
			currentClass = ""
			continue
		}

		if fn, ok := matchTSTopLevelFunc(trimmed, relDir, relPath, lineNum); ok {
			closePrevTSEndLine(functions, lineNum, lastNonEmptyLine)
			functions = append(functions, fn)
			continue
		}

		if fn, ok := tryMatchTSMethod(line, currentClass, relDir, relPath, lineNum); ok {
			closePrevTSEndLine(functions, lineNum, lastNonEmptyLine)
			functions = append(functions, fn)
		}
	}

	// Close the last function's EndLine at end of file.
	closePrevTSEndLineAtEOF(functions, lastNonEmptyLine)

	return functions
}
