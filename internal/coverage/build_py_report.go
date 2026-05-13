//ff:func feature=coverage type=implementation control=iteration dimension=1
//ff:what Builds a Python coverage report from function ranges and coverage data
package coverage

// buildPyReport builds a coverage report from Python function ranges and coverage data.
func buildPyReport(ranges []pyFuncRange, covData *pyCoverageJSON, projectRoot string) *Report {
	report := &Report{AllCovered: true}
	totalLines := 0
	coveredLines := 0

	for _, r := range ranges {
		fc := computePyFuncCoverage(r, covData, projectRoot)
		report.Funcs = append(report.Funcs, fc)
		totalLines += fc.TotalBlocks
		coveredLines += fc.CoveredBlocks
		if fc.CoveredPct < 100 {
			report.AllCovered = false
			appendUncoveredBranches(report, fc)
		}
	}

	if totalLines > 0 {
		report.TotalPct = float64(coveredLines) / float64(totalLines) * 100
	}

	return report
}
