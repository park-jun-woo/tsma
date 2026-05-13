//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs TS test with coverage, parses coverage-final.json, computes per-function coverage
package coverage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Check runs the test with coverage and parses the istanbul/v8 coverage report.
func (c *TSChecker) Check(projectRoot, testFile string, fn *model.Function) (*Report, error) {
	absTest, err := filepath.Abs(testFile)
	if err != nil {
		return nil, fmt.Errorf("resolve test path: %w", err)
	}

	relTest, err := filepath.Rel(projectRoot, absTest)
	if err != nil {
		relTest = absTest
	}

	coverDir := filepath.Join(projectRoot, ".tsma", "coverage")
	if err := os.MkdirAll(coverDir, 0o755); err != nil {
		return nil, fmt.Errorf("create coverage dir: %w", err)
	}

	framework := detectTSCoverageFramework(projectRoot)
	args := buildCoverageArgs(framework, relTest, coverDir)

	cmd := exec.Command("npx", args...)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("test coverage failed: %s\n%s", err, string(output))
	}

	coverageData, err := parseCoverageFinalJSON(coverDir)
	if err != nil {
		return nil, fmt.Errorf("parse coverage report: %w", err)
	}

	ranges := collectTSRanges(fn)

	return buildTSReport(ranges, coverageData, projectRoot), nil
}
