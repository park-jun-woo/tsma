//ff:func feature=endpoint type=implementation control=sequence
//ff:what Walks project Python files to collect FastAPI route registrations
package endpoint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func collectFastapiRoutes(projectRoot string) ([]pyRoute, error) {
	var routes []pyRoute

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if pySkipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".py") {
			return nil
		}
		found, parseErr := parsePyFileForRoutes(path, projectRoot)
		if parseErr != nil {
			return nil
		}
		routes = append(routes, found...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk project: %w", err)
	}

	return routes, nil
}
