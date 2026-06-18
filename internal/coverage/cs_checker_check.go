//ff:func feature=coverage type=implementation control=sequence lang=csharp
//ff:what Runs dotnet test with coverlet and computes per-function coverage from cobertura.xml
package coverage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// Check runs the project's tests with the coverlet XPlat collector (via
// `dotnet test`), collects the Cobertura report under .tsma/coverage, and maps
// it onto the given function's line range. A working .NET SDK plus the
// coverlet.collector package is required. The match is unused: dotnet test
// measures the whole run regardless of the matched file (behavior unchanged).
func (c *CsChecker) Check(projectRoot string, _ match.TestMatch, fn *model.Function) (*Report, error) {
	dotnet, err := findDotnet()
	if err != nil {
		return nil, err
	}

	resultsDir := filepath.Join(projectRoot, ".tsma", "coverage")
	// coverlet's XPlat collector writes each run into a fresh <guid>/ subdir and
	// never overwrites prior reports. findCoberturaReport returns the first match
	// in lexical order, so a stale report from an earlier attempt would be read
	// instead of this run's output. Clear the dir so only the current run's
	// report remains.
	if err := os.RemoveAll(resultsDir); err != nil {
		return nil, fmt.Errorf("clear .tsma coverage dir: %w", err)
	}
	if err := os.MkdirAll(resultsDir, 0o755); err != nil {
		return nil, fmt.Errorf("create .tsma coverage dir: %w", err)
	}

	args := buildCsCoverageArgs(resultsDir)
	if err := runCsCoverage(dotnet, projectRoot, args); err != nil {
		return nil, fmt.Errorf("csharp coverage failed: %w", err)
	}

	reportPath, err := findCoberturaReport(resultsDir)
	if err != nil {
		return nil, fmt.Errorf("locate cobertura xml: %w", err)
	}

	cov, err := parseCobertura(reportPath)
	if err != nil {
		return nil, fmt.Errorf("parse cobertura xml: %w", err)
	}

	ranges := collectCsRanges(fn)

	return buildCsReport(ranges, cov, projectRoot), nil
}
