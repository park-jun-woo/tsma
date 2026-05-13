//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts TS branches within the function range from istanbul coverage data
package coverage

// countTSBranches counts branches within the function range.
func countTSBranches(entry *coverageFinalEntry, r tsFuncRange, fc *FuncCoverage) {
	for branchID, branch := range entry.BranchMap {
		if branch.Loc.Start.Line < r.startLine || branch.Loc.Start.Line > r.endLine {
			continue
		}
		countTSBranchLocations(entry, branchID, branch, r, fc)
	}
}
