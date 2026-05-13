//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs pytest with the --cov flag and JSON output
package coverage

import (
	"fmt"
	"os/exec"
)

// runPytestCov runs pytest with the --cov flag and JSON output.
func runPytestCov(projectRoot, testFile, coverJSONPath string) error {
	cmd := exec.Command("python", "-m", "pytest", testFile,
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
