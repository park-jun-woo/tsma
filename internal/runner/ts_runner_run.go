//ff:func feature=runner type=implementation control=sequence
//ff:what Detects test framework and runs a TS/JS test file via npx
package runner

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
)

// Run executes a TypeScript/JavaScript test file using the detected test framework.
func (r *TSRunner) Run(projectRoot string, m match.TestMatch) (*Result, error) {
	testFile := testFileFromMatch(m)
	absTest, err := filepath.Abs(testFile)
	if err != nil {
		return nil, fmt.Errorf("resolve test path: %w", err)
	}

	relTest, err := filepath.Rel(projectRoot, absTest)
	if err != nil {
		relTest = absTest
	}

	framework := detectTSTestFramework(projectRoot)
	args := buildTestArgs(framework, relTest)

	cmd := exec.Command("npx", args...)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	result := &Result{
		Output: string(output),
		Pass:   err == nil,
	}
	return result, nil
}
