//ff:func feature=coverage type=helper control=sequence
//ff:what Finds the Python binary, preferring python3 over python
package coverage

import "os/exec"

// findCoveragePython returns "python3" if available, otherwise "python".
func findCoveragePython() string {
	if _, err := exec.LookPath("python3"); err == nil {
		return "python3"
	}
	return "python"
}
