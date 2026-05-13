//ff:func feature=endpoint type=implementation control=sequence
//ff:what Finds pages/api/**/*.ts or *.js files and extracts default export handlers
package endpoint

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// detectPagesRouterRoutes scans for Next.js Pages Router API files (pages/api/**/*.ts or *.js).
func detectPagesRouterRoutes(projectRoot string) ([]model.Endpoint, error) {
	var endpoints []model.Endpoint

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if tsSkipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}
		if !isTSOrJSFile(path) {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		if !isUnderDir(relPath, filepath.Join("pages", "api")) {
			return nil
		}

		ep, parseErr := parsePagesRouterFile(path, relPath)
		if parseErr != nil {
			return nil
		}
		if ep != nil {
			endpoints = append(endpoints, *ep)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return endpoints, nil
}
