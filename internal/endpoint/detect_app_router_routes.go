//ff:func feature=endpoint type=implementation control=sequence
//ff:what Finds app/**/route.ts or route.js files and extracts exported HTTP method handlers
package endpoint

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// detectAppRouterRoutes scans for Next.js App Router route files (app/**/route.ts or route.js).
func detectAppRouterRoutes(projectRoot string) ([]model.Endpoint, error) {
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

		base := filepath.Base(path)
		if base != "route.ts" && base != "route.js" {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		if !isUnderDir(relPath, "app") {
			return nil
		}

		eps, parseErr := parseAppRouterFile(path, relPath)
		if parseErr != nil {
			return nil
		}
		endpoints = append(endpoints, eps...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return endpoints, nil
}
