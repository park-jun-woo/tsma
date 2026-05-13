//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts missing lines within the Python function range
package coverage

// countPyMissingLines counts missing lines within the function range.
func countPyMissingLines(fileCov *pyCoverageFile, r pyFuncRange, fc *FuncCoverage) {
	for _, line := range fileCov.MissingLines {
		if line >= r.startLine && line <= r.endLine {
			fc.TotalBlocks++
			fc.UncoveredLines = append(fc.UncoveredLines, line)
		}
	}
}
