//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs coverage.py as a fallback to collect Python test coverage
package coverage

import (
	"fmt"
	"os/exec"
)

// runCoveragePy runs coverage.py as a fallback.
func runCoveragePy(projectRoot, testFile, coverJSONPath string) error {
	runCmd := exec.Command("python", "-m", "coverage", "run", "-m", "pytest", testFile)
	runCmd.Dir = projectRoot

	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("coverage run: %s\n%s", err, string(output))
	}

	jsonCmd := exec.Command("python", "-m", "coverage", "json", "-o", coverJSONPath)
	jsonCmd.Dir = projectRoot

	output, err = jsonCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("coverage json: %s\n%s", err, string(output))
	}

	return nil
}
