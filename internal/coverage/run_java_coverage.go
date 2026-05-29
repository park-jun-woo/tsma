//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs the Java build tool to execute tests and generate a JaCoCo report
package coverage

import (
	"fmt"
	"os/exec"
)

// runJavaCoverage runs the build tool (mvn/gradle) with the coverage arguments
// in projectRoot, producing a JaCoCo XML report. A working JDK + build tool +
// JaCoCo plugin is required (E2E only; not exercised in sandbox environments).
func runJavaCoverage(bin, projectRoot string, args []string) error {
	cmd := exec.Command(bin, args...)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("java coverage run: %s\n%s", err, string(output))
	}
	return nil
}
