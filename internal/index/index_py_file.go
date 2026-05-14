//ff:func feature=index type=implementation control=iteration dimension=1
//ff:what Parses a single Python file and extracts function declarations line by line
package index

import (
	"bufio"
	"os"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// indexPyFile parses a single Python file and extracts function declarations.
func indexPyFile(relPath, absPath string) []model.Function {
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

		isNonEmpty := strings.TrimSpace(line) != ""

		if m := pyClassDefPattern.FindStringSubmatch(line); m != nil {
			if isNonEmpty {
				lastNonEmptyLine = lineNum
			}
			currentClass = m[2]
			classIndent = pyIndent(m[1])
			continue
		}

		fn, newClass := matchPyFunc(line, lineNum, relPath, relDir, currentClass, classIndent)
		if fn != nil {
			// Close the previous function's EndLine before appending the new one.
			// lastNonEmptyLine holds the last non-empty line BEFORE this def line.
			if n := len(functions); n > 0 && functions[n-1].EndLine == 0 {
				functions[n-1].EndLine = lastNonEmptyLine
			}
			functions = append(functions, *fn)
		}
		currentClass = newClass

		if isNonEmpty {
			lastNonEmptyLine = lineNum
		}
	}

	// Close the last function's EndLine at end of file.
	if n := len(functions); n > 0 && functions[n-1].EndLine == 0 {
		functions[n-1].EndLine = lastNonEmptyLine
	}

	return functions
}
