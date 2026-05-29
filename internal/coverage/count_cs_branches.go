//ff:func feature=coverage type=helper control=iteration dimension=1 lang=csharp
//ff:what Counts branch blocks within the C# function range and records uncovered branch lines
package coverage

// countCsBranches counts branch outcomes on each branch line within the function
// range. Cobertura reports branch lines with a condition-coverage attribute of
// the form "covered% (covered/total)"; each branch outcome is one block, and any
// missed outcome records the line as uncovered.
func countCsBranches(fileCov *csFile, r csFuncRange, fc *FuncCoverage) {
	for _, ln := range fileCov.Lines {
		if ln.Number < r.startLine || ln.Number > r.endLine {
			continue
		}
		if ln.Branch != "true" {
			continue
		}
		covered, total := parseCsConditionCoverage(ln.ConditionCoverage)
		if total == 0 {
			continue
		}
		fc.TotalBlocks += total
		fc.CoveredBlocks += covered
		for i := 0; i < total-covered; i++ {
			fc.UncoveredLines = append(fc.UncoveredLines, ln.Number)
		}
	}
}
