//ff:func feature=coverage type=implementation control=iteration dimension=1 lang=go
//ff:what Runs go test once for a package then computes per-function coverage reports for all its functions
package coverage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/runner"
)

// CheckPkg measures coverage for every function in a single package with one
// `go test -coverprofile` run. pkgDir is the absolute directory of the package;
// funcs are the functions in that package; testFuncs is the union of test
// function names attributed to those functions (used as the anchored -run
// filter). The profile is parsed once and reused per function via
// collectRanges+buildGoReport, so the whole package is measured in one test run
// (the efficiency core of the batch first scan). The returned map is keyed by
// the same *model.Function pointers passed in. A test run / compile failure
// returns an error so the caller can skip this package without aborting the scan.
func (c *GoChecker) CheckPkg(projectRoot, pkgDir string, funcs []*model.Function, testFuncs []string) (map[*model.Function]*Report, error) {
	if len(funcs) == 0 {
		return map[*model.Function]*Report{}, nil
	}

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
	if pattern := runner.AnchorRunPattern(testFuncs); pattern != "" {
		args = append(args, "-run", pattern)
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

	reports := make(map[*model.Function]*Report, len(funcs))
	for _, fn := range funcs {
		ranges := collectRanges(fn)
		reports[fn] = buildGoReport(ranges, blocks, projectRoot)
	}
	return reports, nil
}
