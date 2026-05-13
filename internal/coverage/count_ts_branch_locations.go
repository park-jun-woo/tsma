//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts individual TS branch locations within a function range
package coverage

// countTSBranchLocations counts individual branch locations within range.
func countTSBranchLocations(entry *coverageFinalEntry, branchID string, branch coverageBranch, r tsFuncRange, fc *FuncCoverage) {
	counts := entry.B[branchID]
	for i, loc := range branch.Locations {
		if loc.Start.Line < r.startLine || loc.Start.Line > r.endLine {
			continue
		}
		fc.TotalBlocks++
		if i < len(counts) && counts[i] > 0 {
			fc.CoveredBlocks++
		} else {
			fc.UncoveredLines = append(fc.UncoveredLines, loc.Start.Line)
		}
	}
}
