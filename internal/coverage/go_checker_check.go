//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs go test with coverprofile and computes per-function coverage
package coverage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
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

	cmd := exec.Command("go", "test", "-count=1",
		"-coverprofile="+coverFile,
		"-covermode=set",
		pkgPath,
	)
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
