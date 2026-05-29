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
			updateLastNonEmpty(isNonEmpty, lineNum, &lastNonEmptyLine)
			currentClass = m[2]
			classIndent = pyIndent(m[1])
			continue
		}

		fn, newClass := matchPyFunc(line, lineNum, relPath, relDir, currentClass, classIndent)
		if fn != nil {
			closePrevEndLine(functions, lastNonEmptyLine)
			functions = append(functions, *fn)
		}
		currentClass = newClass

		updateLastNonEmpty(isNonEmpty, lineNum, &lastNonEmptyLine)
	}

	// Close the last function's EndLine at end of file.
	closePrevEndLine(functions, lastNonEmptyLine)

	return functions
}
