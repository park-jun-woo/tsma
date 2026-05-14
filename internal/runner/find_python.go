//ff:func feature=runner type=helper control=sequence
//ff:what Finds the Python binary, preferring python3 over python
package runner

import "os/exec"

// findPython returns "python3" if available, otherwise "python".
func findPython() string {
	if _, err := exec.LookPath("python3"); err == nil {
		return "python3"
	}
	return "python"
}
