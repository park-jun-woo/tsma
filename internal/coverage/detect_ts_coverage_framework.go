//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Reads package.json to detect vitest or jest for coverage command selection
package coverage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// detectTSCoverageFramework reads package.json to determine the test framework.
func detectTSCoverageFramework(projectRoot string) tsCoverageFramework {
	pkgPath := filepath.Join(projectRoot, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return coverVitest
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return coverVitest
	}

	allDeps := make(map[string]bool)
	for k := range pkg.DevDependencies {
		allDeps[k] = true
	}
	for k := range pkg.Dependencies {
		allDeps[k] = true
	}

	if allDeps["vitest"] {
		return coverVitest
	}
	if allDeps["jest"] {
		return coverJest
	}
	return coverVitest
}
