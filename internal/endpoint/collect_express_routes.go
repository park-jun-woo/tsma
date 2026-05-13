//ff:func feature=endpoint type=implementation control=sequence
//ff:what Walks project TS/JS files to collect Express route registrations
package endpoint

import (
	"fmt"
	"os"
	"path/filepath"
)

func collectExpressRoutes(projectRoot string) ([]tsRouteRegistration, error) {
	var regs []tsRouteRegistration

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
		found, parseErr := parseExpressRoutes(path, projectRoot)
		if parseErr != nil {
			return nil
		}
		regs = append(regs, found...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk project: %w", err)
	}

	return regs, nil
}
