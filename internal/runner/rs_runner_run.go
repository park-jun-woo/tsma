//ff:func feature=runner type=implementation control=sequence
//ff:what Executes Rust tests for the given test file using cargo test
package runner

import (
	"os/exec"

	"github.com/park-jun-woo/tsma/internal/match"
)

// E2E note: running this path requires a working `cargo` toolchain, which is
// not available in CI/sandbox environments. The argument builder is split out
// (buildCargoTestArgs) for environment-independent unit testing.

// Run executes the given Rust test file against the project using cargo test.
func (r *RsRunner) Run(projectRoot string, m match.TestMatch) (*Result, error) {
	testFile := testFileFromMatch(m)
	cargo, err := findCargo()
	if err != nil {
		return nil, err
	}

	args := buildCargoTestArgs(testFile)
	cmd := exec.Command(cargo, args...)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	result := &Result{
		Output: string(output),
		Pass:   err == nil,
	}
	return result, nil
}
