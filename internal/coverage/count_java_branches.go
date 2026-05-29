//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts branch blocks within the Java function range and records uncovered branch lines
package coverage

// countJavaBranches counts branch outcomes on each line within the function
// range. JaCoCo reports per-line covered (Cb) and missed (Mb) branch counts;
// each branch outcome is one block, and missed outcomes record the line as
// uncovered.
func countJavaBranches(fileCov *jacocoFile, r javaFuncRange, fc *FuncCoverage) {
	for _, ln := range fileCov.Lines {
		if ln.Nr < r.startLine || ln.Nr > r.endLine {
			continue
		}
		fc.TotalBlocks += ln.Cb + ln.Mb
		fc.CoveredBlocks += ln.Cb
		for i := 0; i < ln.Mb; i++ {
			fc.UncoveredLines = append(fc.UncoveredLines, ln.Nr)
		}
	}
}
