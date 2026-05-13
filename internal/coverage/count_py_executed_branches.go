//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts executed branches within the Python function range
package coverage

// countPyExecutedBranches counts executed branches within the function range.
func countPyExecutedBranches(fileCov *pyCoverageFile, r pyFuncRange, fc *FuncCoverage) {
	for _, branch := range fileCov.ExecutedBranches {
		if len(branch) >= 2 && branch[0] >= r.startLine && branch[0] <= r.endLine {
			fc.TotalBlocks++
			fc.CoveredBlocks++
		}
	}
}
