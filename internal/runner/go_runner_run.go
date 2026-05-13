//ff:func feature=runner type=implementation control=sequence
//ff:what Executes a Go test file by resolving its package and running go test
package runner

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// Run executes the test file against the original code.
func (r *GoRunner) Run(projectRoot, testFile string) (*Result, error) {
	absTest, err := filepath.Abs(testFile)
	if err != nil {
		return nil, fmt.Errorf("resolve test path: %w", err)
	}

	pkgPath, err := resolveGoPkgPath(projectRoot, absTest)
	if err != nil {
		return nil, err
	}

	testFuncs, err := extractTestFuncs(absTest)
	if err != nil {
		return nil, fmt.Errorf("extract test functions: %w", err)
	}

	args := buildGoTestArgs(pkgPath, testFuncs)

	cmd := exec.Command("go", args...)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	result := &Result{
		Output: string(output),
		Pass:   err == nil,
	}
	return result, nil
}
