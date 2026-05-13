//ff:func feature=coverage type=implementation control=iteration dimension=1
//ff:what Builds a Go coverage report from function ranges and coverage blocks
package coverage

// buildGoReport builds a coverage report from function ranges and coverage blocks.
func buildGoReport(ranges []funcRange, blocks []coverBlock, projectRoot string) *Report {
	report := &Report{AllCovered: true}
	totalBlocks := 0
	coveredBlocks := 0

	for _, r := range ranges {
		fc := computeFuncCoverage(r, blocks, projectRoot)
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
