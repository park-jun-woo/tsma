//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs coverage.py as a fallback to collect Python test coverage
package coverage

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// runCoveragePy runs coverage.py as a fallback using unittest.
func runCoveragePy(projectRoot, testFile, coverJSONPath string) error {
	python := findCoveragePython()

	absTest := filepath.Join(projectRoot, testFile)
	rel, _ := filepath.Rel(projectRoot, absTest)
	modulePath := strings.TrimSuffix(rel, ".py")
	modulePath = strings.ReplaceAll(modulePath, string(filepath.Separator), ".")

	runCmd := exec.Command(python, "-m", "coverage", "run", "-m", "unittest", modulePath, "-v")
	runCmd.Dir = projectRoot

	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("coverage run: %s\n%s", err, string(output))
	}

	jsonCmd := exec.Command(python, "-m", "coverage", "json", "-o", coverJSONPath)
	jsonCmd.Dir = projectRoot

	output, err = jsonCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("coverage json: %s\n%s", err, string(output))
	}

	return nil
}
