//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts executed lines within the Python function range
package coverage

// countPyExecutedLines counts executed lines within the function range.
func countPyExecutedLines(fileCov *pyCoverageFile, r pyFuncRange, fc *FuncCoverage) {
	for _, line := range fileCov.ExecutedLines {
		if line >= r.startLine && line <= r.endLine {
			fc.TotalBlocks++
			fc.CoveredBlocks++
		}
	}
}
