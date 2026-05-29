//ff:func feature=coverage type=implementation control=iteration dimension=1
//ff:what Builds a Rust coverage report from function ranges and llvm-cov data
package coverage

// buildRsReport builds a coverage report from Rust function ranges and llvm-cov data.
func buildRsReport(ranges []rsFuncRange, cov *llvmCovJSON, projectRoot string) *Report {
	report := &Report{AllCovered: true}
	totalBlocks := 0
	coveredBlocks := 0

	for _, r := range ranges {
		fc := computeRsFuncCoverage(r, cov, projectRoot)
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
