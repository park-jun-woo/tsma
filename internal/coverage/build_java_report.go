//ff:func feature=coverage type=implementation control=iteration dimension=1
//ff:what Builds a Java coverage report from function ranges and JaCoCo coverage data
package coverage

// buildJavaReport builds a coverage report from Java function ranges and JaCoCo
// coverage data.
func buildJavaReport(ranges []javaFuncRange, cov *jacocoCoverage, projectRoot string) *Report {
	report := &Report{AllCovered: true}
	totalBlocks := 0
	coveredBlocks := 0

	for _, r := range ranges {
		fc := computeJavaFuncCoverage(r, cov, projectRoot)
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
