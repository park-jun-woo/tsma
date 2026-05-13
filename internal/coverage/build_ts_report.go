//ff:func feature=coverage type=implementation control=iteration dimension=1
//ff:what Builds a TS coverage report from function ranges and istanbul coverage data
package coverage

// buildTSReport builds a coverage report from TS function ranges and coverage data.
func buildTSReport(ranges []tsFuncRange, coverageData map[string]coverageFinalEntry, projectRoot string) *Report {
	report := &Report{AllCovered: true}
	totalStmts := 0
	coveredStmts := 0

	for _, r := range ranges {
		fc := computeTSFuncCoverage(r, coverageData, projectRoot)
		report.Funcs = append(report.Funcs, fc)
		totalStmts += fc.TotalBlocks
		coveredStmts += fc.CoveredBlocks
		if fc.CoveredPct < 100 {
			report.AllCovered = false
			appendUncoveredBranches(report, fc)
		}
	}

	if totalStmts > 0 {
		report.TotalPct = float64(coveredStmts) / float64(totalStmts) * 100
	}

	return report
}
