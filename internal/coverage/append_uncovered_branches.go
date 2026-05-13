//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Appends uncovered lines from a FuncCoverage to the report's Uncovered list
package coverage

// appendUncoveredBranches adds uncovered lines from a FuncCoverage to the report.
func appendUncoveredBranches(report *Report, fc FuncCoverage) {
	for _, line := range fc.UncoveredLines {
		report.Uncovered = append(report.Uncovered, UncoveredBranch{
			File: fc.File,
			Line: line,
		})
	}
}
