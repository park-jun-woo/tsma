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
//
// SSOT note (Phase006/D2): this stage assumes pytest unconditionally and is
// already correct for pytest projects, so its execution path is unchanged. The
// single source of truth for "is this a pytest project?" is
// detect.DetectPytest, which the runner stage now also uses. Should a future
// non-pytest coverage branch be added here, it must gate on detect.DetectPytest
// so the runner and coverage stages stay in agreement (the asymmetry behind
// BUG-001).
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
