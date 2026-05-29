//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs go test with coverprofile filtered by the match's anchored test funcs and computes per-function coverage
package coverage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/runner"
)

// Check runs go test with coverage and parses the profile. The package path is
// derived from the match's files (assumed same package) and the -run filter is
// the anchored union of the match's test functions, so only the attributed
// tests run.
func (c *GoChecker) Check(projectRoot string, m match.TestMatch, fn *model.Function) (*Report, error) {
	if len(m.Files) == 0 {
		return nil, fmt.Errorf("no test files in match")
	}

	absTest, err := filepath.Abs(m.Files[0])
	if err != nil {
		return nil, fmt.Errorf("resolve test path: %w", err)
	}

	pkgDir := filepath.Dir(absTest)
	relPkg, err := filepath.Rel(projectRoot, pkgDir)
	if err != nil {
		return nil, fmt.Errorf("resolve relative package: %w", err)
	}
	pkgPath := "./" + filepath.ToSlash(relPkg)

	coverFile := filepath.Join(projectRoot, ".tsma", "cover.out")
	os.MkdirAll(filepath.Dir(coverFile), 0o755)

	args := []string{"test", "-count=1",
		"-coverprofile=" + coverFile,
		"-covermode=set",
	}

	testFuncs, err := runner.ResolveTestFuncs(m)
	if err == nil {
		if pattern := runner.AnchorRunPattern(testFuncs); pattern != "" {
			args = append(args, "-run", pattern)
		}
	}

	args = append(args, pkgPath)

	cmd := exec.Command("go", args...)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("go test coverage failed: %s\n%s", err, string(output))
	}

	blocks, err := parseCoverProfile(coverFile)
	if err != nil {
		return nil, fmt.Errorf("parse coverage profile: %w", err)
	}

	ranges := collectRanges(fn)

	return buildGoReport(ranges, blocks, projectRoot), nil
}
