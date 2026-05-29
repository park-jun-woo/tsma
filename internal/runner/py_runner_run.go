//ff:func feature=runner type=implementation control=sequence
//ff:what Executes a Python test file using pytest or unittest
package runner

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/tsma/internal/match"
)

// Run executes the given Python test file against the project.
func (r *PyRunner) Run(projectRoot string, m match.TestMatch) (*Result, error) {
	testFile := testFileFromMatch(m)
	absTest := filepath.Join(projectRoot, testFile)

	usePytest := detectPytest(projectRoot)
	python := findPython()

	var cmd *exec.Cmd
	if usePytest {
		cmd = exec.Command(python, "-m", "pytest", absTest, "-v")
	} else {
		rel, _ := filepath.Rel(projectRoot, absTest)
		modulePath := strings.TrimSuffix(rel, ".py")
		modulePath = strings.ReplaceAll(modulePath, string(filepath.Separator), ".")
		cmd = exec.Command(python, "-m", "unittest", modulePath, "-v")
	}
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	result := &Result{
		Output: string(output),
		Pass:   err == nil,
	}
	return result, nil
}
