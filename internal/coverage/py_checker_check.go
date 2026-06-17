//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs Python tests with coverage and computes per-function coverage
package coverage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// Check runs pytest with coverage and parses the results.
func (c *PyChecker) Check(projectRoot string, m match.TestMatch, fn *model.Function) (*Report, error) {
	testFile := testFileFromMatch(m)
	// Validate the path is resolvable (e.g. cwd still exists), but pass the
	// ROOT-RELATIVE testFile to the run helpers below: they Join it with
	// projectRoot themselves, so handing them an absolute path would double the
	// prefix (Join("/root", "/root/x") = "/root/root/x"). This is the seam the
	// loop's .tsma/test backing flows through.
	if _, err := filepath.Abs(testFile); err != nil {
		return nil, fmt.Errorf("resolve test path: %w", err)
	}

	tmDir := filepath.Join(projectRoot, ".tsma")
	if err := os.MkdirAll(tmDir, 0o755); err != nil {
		return nil, fmt.Errorf("create .tsma dir: %w", err)
	}

	coverJSONPath := filepath.Join(tmDir, "coverage.json")

	err := runPytestCov(projectRoot, testFile, coverJSONPath)
	if err != nil {
		err = runCoveragePy(projectRoot, testFile, coverJSONPath)
		if err != nil {
			return nil, fmt.Errorf("python coverage failed: %w", err)
		}
	}

	covData, err := parsePyCoverageJSON(coverJSONPath)
	if err != nil {
		return nil, fmt.Errorf("parse coverage json: %w", err)
	}

	ranges := collectPyRanges(fn)

	return buildPyReport(ranges, covData, projectRoot), nil
}
