//ff:func feature=endpoint type=implementation control=iteration dimension=1
//ff:what Looks for a Django view definition in a single Python file
package endpoint

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// findDjangoViewInFile looks for the view definition in a single .py file.
func findDjangoViewInFile(filePath, projectRoot string, route djangoRoute) []model.Endpoint {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	relPath, _ := filepath.Rel(projectRoot, filePath)

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil
	}

	for i, line := range lines {
		lineNum := i + 1

		if eps := matchFuncView(line, lines, i, lineNum, relPath, route); eps != nil {
			return eps
		}

		if eps := matchClassView(line, lines, i, lineNum, relPath, route); eps != nil {
			return eps
		}
	}

	return nil
}
