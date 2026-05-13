//ff:func feature=runner type=implementation control=sequence
//ff:what Executes a Python test file using pytest or unittest
package runner

import (
	"os/exec"
	"path/filepath"
)

// Run executes the given Python test file against the project.
func (r *PyRunner) Run(projectRoot, testFile string) (*Result, error) {
	absTest, err := filepath.Abs(testFile)
	if err != nil {
		return &Result{Pass: false, Output: "failed to resolve test path: " + err.Error()}, nil
	}

	usePytest := detectPytest(projectRoot)

	var cmd *exec.Cmd
	if usePytest {
		cmd = exec.Command("python", "-m", "pytest", absTest, "-v")
	} else {
		cmd = exec.Command("python", "-m", "unittest", absTest, "-v")
	}
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	result := &Result{
		Output: string(output),
		Pass:   err == nil,
	}
	return result, nil
}
