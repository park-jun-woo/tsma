//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts executable line blocks within the Java function range and records uncovered lines
package coverage

// countJavaLines counts executable lines (those carrying instructions) whose
// number lies within the function range. Each such line is one block; it is
// covered when covered-instructions (Ci) is positive, otherwise its number is
// recorded as uncovered.
func countJavaLines(fileCov *jacocoFile, r javaFuncRange, fc *FuncCoverage) {
	for _, ln := range fileCov.Lines {
		if ln.Nr < r.startLine || ln.Nr > r.endLine {
			continue
		}
		if ln.Mi == 0 && ln.Ci == 0 {
			continue // no instructions on this line
		}
		fc.TotalBlocks++
		if ln.Ci > 0 {
			fc.CoveredBlocks++
		} else {
			fc.UncoveredLines = append(fc.UncoveredLines, ln.Nr)
		}
	}
}
