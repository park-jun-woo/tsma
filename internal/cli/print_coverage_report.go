//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints coverage details for each function in the report
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/coverage"
)

// printCoverageReport prints coverage details for each function.
func printCoverageReport(report *coverage.Report) {
	for _, fc := range report.Funcs {
		pct := fmt.Sprintf("%.0f%%", fc.CoveredPct)
		blocks := fmt.Sprintf("(%d/%d blocks)", fc.CoveredBlocks, fc.TotalBlocks)
		loc := fmt.Sprintf("%s:%d-%d", fc.File, fc.StartLine, fc.EndLine)
		fmt.Printf("  %-40s %s %s\n", loc, pct, blocks)
		printUncoveredLines(fc.UncoveredLines)
	}
}
