//ff:func feature=endpoint type=implementation control=sequence
//ff:what Searches project files for the Django view definition and returns endpoints
package endpoint

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// resolveDjangoView searches the project for the view definition and returns endpoints.
func resolveDjangoView(projectRoot string, route djangoRoute) []model.Endpoint {
	var endpoints []model.Endpoint

	filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if djangoSkipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".py") {
			return nil
		}

		found := findDjangoViewInFile(path, projectRoot, route)
		endpoints = append(endpoints, found...)

		if len(endpoints) > 0 {
			return filepath.SkipAll
		}
		return nil
	})

	if len(endpoints) == 0 {
		endpoints = append(endpoints, model.Endpoint{
			Name:   route.viewName,
			Method: "ANY",
			Path:   route.path,
			Status: model.StatusTodo,
		})
	}

	return endpoints
}
