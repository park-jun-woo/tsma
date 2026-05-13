//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts missing branches within the Python function range
package coverage

// countPyMissingBranches counts missing branches within the function range.
func countPyMissingBranches(fileCov *pyCoverageFile, r pyFuncRange, fc *FuncCoverage) {
	for _, branch := range fileCov.MissingBranches {
		if len(branch) >= 2 && branch[0] >= r.startLine && branch[0] <= r.endLine {
			fc.TotalBlocks++
			fc.UncoveredLines = append(fc.UncoveredLines, branch[0])
		}
	}
}
