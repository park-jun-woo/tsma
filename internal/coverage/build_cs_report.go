//ff:func feature=coverage type=implementation control=iteration dimension=1 lang=csharp
//ff:what Builds a C# coverage report from function ranges and Cobertura coverage data
package coverage

// buildCsReport builds a coverage report from C# function ranges and Cobertura
// coverage data.
func buildCsReport(ranges []csFuncRange, cov *csCoverage, projectRoot string) *Report {
	report := &Report{AllCovered: true}
	totalBlocks := 0
	coveredBlocks := 0

	for _, r := range ranges {
		fc := computeCsFuncCoverage(r, cov, projectRoot)
		report.Funcs = append(report.Funcs, fc)
		totalBlocks += fc.TotalBlocks
		coveredBlocks += fc.CoveredBlocks
		if fc.CoveredPct < 100 {
			report.AllCovered = false
			appendUncoveredBranches(report, fc)
		}
	}

	if totalBlocks > 0 {
		report.TotalPct = float64(coveredBlocks) / float64(totalBlocks) * 100
	}

	return report
}
