//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs pytest with the --cov flag and JSON output
package coverage

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// runPytestCov runs pytest with the --cov flag and JSON output.
func runPytestCov(projectRoot, testFile, coverJSONPath string) error {
	python := findCoveragePython()
	absTest := filepath.Join(projectRoot, testFile)

	cmd := exec.Command(python, "-m", "pytest", absTest,
		"--cov",
		"--cov-report=json:"+coverJSONPath,
		"--cov-branch",
	)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pytest-cov: %s\n%s", err, string(output))
	}
	return nil
}
