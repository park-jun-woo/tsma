//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs go test with coverprofile filtered by test file functions and computes per-function coverage
package coverage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/runner"
)

// Check runs go test with coverage and parses the profile.
func (c *GoChecker) Check(projectRoot, testFile string, fn *model.Function) (*Report, error) {
	absTest, err := filepath.Abs(testFile)
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

	testFuncs, err := runner.ExtractTestFuncs(absTest)
	if err == nil && len(testFuncs) > 0 {
		args = append(args, "-run", strings.Join(testFuncs, "|"))
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
