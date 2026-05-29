//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts branch sides within the Rust function range and records uncovered branch lines
package coverage

// countRsBranches counts both sides (true/false) of each branch region whose
// start line lies within the function range. A side is covered when its
// execution count is positive.
func countRsBranches(fileCov *llvmCovFile, r rsFuncRange, fc *FuncCoverage) {
	for _, br := range fileCov.Branches {
		if br.LineStart < r.startLine || br.LineStart > r.endLine {
			continue
		}
		// True side.
		fc.TotalBlocks++
		if br.ExecCount > 0 {
			fc.CoveredBlocks++
		} else {
			fc.UncoveredLines = append(fc.UncoveredLines, br.LineStart)
		}
		// False side.
		fc.TotalBlocks++
		if br.FalseExecCount > 0 {
			fc.CoveredBlocks++
		} else {
			fc.UncoveredLines = append(fc.UncoveredLines, br.LineStart)
		}
	}
}
