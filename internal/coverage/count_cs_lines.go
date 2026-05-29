//ff:func feature=coverage type=helper control=iteration dimension=1 lang=csharp
//ff:what Counts executable line blocks within the C# function range and records uncovered lines
package coverage

// countCsLines counts executable lines whose number lies within the function
// range. Each Cobertura <line> is one block; it is covered when hits is
// positive, otherwise its number is recorded as uncovered.
func countCsLines(fileCov *csFile, r csFuncRange, fc *FuncCoverage) {
	for _, ln := range fileCov.Lines {
		if ln.Number < r.startLine || ln.Number > r.endLine {
			continue
		}
		fc.TotalBlocks++
		if ln.Hits > 0 {
			fc.CoveredBlocks++
		} else {
			fc.UncoveredLines = append(fc.UncoveredLines, ln.Number)
		}
	}
}
