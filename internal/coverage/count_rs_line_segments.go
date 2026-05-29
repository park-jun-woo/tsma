//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Counts executable line segments within the Rust function range and records uncovered lines
package coverage

// countRsLineSegments counts region-entry line segments inside the function
// range. Each such segment is one block; it is covered when its execution count
// is positive, otherwise its line is recorded as uncovered.
func countRsLineSegments(fileCov *llvmCovFile, r rsFuncRange, fc *FuncCoverage) {
	for _, seg := range fileCov.Segments {
		if !seg.IsRegionEntry || !seg.HasCount {
			continue
		}
		if seg.Line < r.startLine || seg.Line > r.endLine {
			continue
		}
		fc.TotalBlocks++
		if seg.Count > 0 {
			fc.CoveredBlocks++
		} else {
			fc.UncoveredLines = append(fc.UncoveredLines, seg.Line)
		}
	}
}
