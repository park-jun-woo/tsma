//ff:func feature=endpoint type=implementation control=iteration dimension=1
//ff:what Scans a single Python file for FastAPI route decorators and associated def lines
package endpoint

import (
	"bufio"
	"os"
	"path/filepath"
)

// parsePyFileForRoutes scans a single Python file for FastAPI route decorators.
func parsePyFileForRoutes(filePath, projectRoot string) ([]pyRoute, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	relPath, _ := filepath.Rel(projectRoot, filePath)

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return matchPyRoutes(lines, relPath), nil
}
