//ff:func feature=endpoint type=implementation control=sequence
//ff:what Walks project Go files to collect Chi route registrations
package endpoint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func collectChiRoutes(projectRoot string) ([]routeRegistration, error) {
	var regs []routeRegistration

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if base == "vendor" || base == ".git" || base == ".tsma" || base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		found, parseErr := parseFileForChiRoutes(path, projectRoot)
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
