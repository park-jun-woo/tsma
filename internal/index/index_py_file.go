//ff:func feature=index type=implementation control=iteration dimension=1
//ff:what Parses a single Python file and extracts function declarations line by line
package index

import (
	"bufio"
	"os"

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

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if m := pyClassDefPattern.FindStringSubmatch(line); m != nil {
			currentClass = m[2]
			classIndent = pyIndent(m[1])
			continue
		}

		fn, newClass := matchPyFunc(line, lineNum, relPath, relDir, currentClass, classIndent)
		if fn != nil {
			functions = append(functions, *fn)
		}
		currentClass = newClass
	}

	return functions
}
