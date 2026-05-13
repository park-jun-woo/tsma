//ff:func feature=runner type=helper control=iteration dimension=1
//ff:what Reads package.json dependencies to detect vitest, jest, or mocha
package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// detectTSTestFramework reads package.json to determine which test framework is used.
func detectTSTestFramework(projectRoot string) tsTestFramework {
	pkgPath := filepath.Join(projectRoot, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return frameworkVitest // fallback
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return frameworkVitest // fallback
	}

	allDeps := make(map[string]bool)
	for k := range pkg.DevDependencies {
		allDeps[k] = true
	}
	for k := range pkg.Dependencies {
		allDeps[k] = true
	}

	if allDeps["vitest"] {
		return frameworkVitest
	}
	if allDeps["jest"] {
		return frameworkJest
	}
	if allDeps["mocha"] {
		return frameworkMocha
	}

	return frameworkVitest // fallback
}
