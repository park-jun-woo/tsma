//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs Python tests with coverage and computes per-function coverage
package coverage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Check runs pytest with coverage and parses the results.
func (c *PyChecker) Check(projectRoot, testFile string, fn *model.Function) (*Report, error) {
	absTest, err := filepath.Abs(testFile)
	if err != nil {
		return nil, fmt.Errorf("resolve test path: %w", err)
	}

	tmDir := filepath.Join(projectRoot, ".tsma")
	if err := os.MkdirAll(tmDir, 0o755); err != nil {
		return nil, fmt.Errorf("create .tsma dir: %w", err)
	}

	coverJSONPath := filepath.Join(tmDir, "coverage.json")

	err = runPytestCov(projectRoot, absTest, coverJSONPath)
	if err != nil {
		err = runCoveragePy(projectRoot, absTest, coverJSONPath)
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
