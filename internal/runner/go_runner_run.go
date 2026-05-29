//ff:func feature=runner type=implementation control=sequence
//ff:what Executes Go tests for a content-aware match, resolving the package and anchored -run set
package runner

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
)

// Run executes the matched tests against the original code. It resolves the
// package from the match's files (assumed same package) and runs the union of
// the match's test functions with an anchored -run filter. When the match
// carries no explicit test functions (nil TestFuncs), it extracts them from the
// matched files, preserving the legacy single-file behavior.
func (r *GoRunner) Run(projectRoot string, m match.TestMatch) (*Result, error) {
	if len(m.Files) == 0 {
		return nil, fmt.Errorf("no test files in match")
	}

	absFirst, err := filepath.Abs(m.Files[0])
	if err != nil {
		return nil, fmt.Errorf("resolve test path: %w", err)
	}

	pkgPath, err := resolveGoPkgPath(projectRoot, absFirst)
	if err != nil {
		return nil, err
	}

	testFuncs, err := ResolveTestFuncs(m)
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
