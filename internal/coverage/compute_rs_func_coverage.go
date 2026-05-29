//ff:func feature=coverage type=implementation control=sequence
//ff:what Computes Rust function coverage by mapping llvm-cov segments and branches onto a function range
package coverage

import "fmt"

// computeRsFuncCoverage computes coverage for a single Rust function range,
// counting executable line segments and branch sides within the range.
func computeRsFuncCoverage(r rsFuncRange, cov *llvmCovJSON, projectRoot string) FuncCoverage {
	fc := FuncCoverage{
		File:      r.file,
		StartLine: r.startLine,
		EndLine:   r.endLine,
		Key:       fmt.Sprintf("%s:%d-%d", r.file, r.startLine, r.endLine),
	}

	if cov == nil {
		fc.CoveredPct = 100
		return fc
	}

	fileCov := findRsCoverageFile(cov, r.file, projectRoot)
	if fileCov == nil {
		fc.CoveredPct = 100
		return fc
	}

	countRsLineSegments(fileCov, r, &fc)
	countRsBranches(fileCov, r, &fc)

	if fc.TotalBlocks > 0 {
		fc.CoveredPct = float64(fc.CoveredBlocks) / float64(fc.TotalBlocks) * 100
	} else {
		fc.CoveredPct = 100
	}

	return fc
}
