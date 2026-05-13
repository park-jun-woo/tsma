//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts TS statements within the function range from istanbul coverage data
package coverage

// countTSStatements counts statements within the function range.
func countTSStatements(entry *coverageFinalEntry, r tsFuncRange, fc *FuncCoverage) {
	for stmtID, stmtRange := range entry.StatementMap {
		if stmtRange.Start.Line < r.startLine || stmtRange.Start.Line > r.endLine {
			continue
		}
		fc.TotalBlocks++
		count := entry.S[stmtID]
		if count > 0 {
			fc.CoveredBlocks++
		} else {
			fc.UncoveredLines = append(fc.UncoveredLines, stmtRange.Start.Line)
		}
	}
}
