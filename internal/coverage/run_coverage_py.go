//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs coverage.py as a fallback to collect Python test coverage (branch mode, pytest by path)
package coverage

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// runCoveragePy runs coverage.py as a fallback when pytest-cov is unavailable.
// It drives pytest by file PATH under `coverage run --branch` (Phase005b §3:
// branch mode is explicit), which — unlike a unittest module path — collects a
// test file anywhere on disk, including the loop's .tsma/test backing whose
// directory is not an importable package.
func runCoveragePy(projectRoot, testFile, coverJSONPath string) error {
	python := findCoveragePython()

	absTest := filepath.Join(projectRoot, testFile)

	runCmd := exec.Command(python, "-m", "coverage", "run", "--branch", "-m", "pytest", absTest)
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
