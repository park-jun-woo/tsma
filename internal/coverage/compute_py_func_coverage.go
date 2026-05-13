//ff:func feature=coverage type=implementation control=sequence
//ff:what Computes Python function coverage by matching coverage data against a function range
package coverage

import "fmt"

// computePyFuncCoverage computes coverage for a single function range.
func computePyFuncCoverage(r pyFuncRange, covData *pyCoverageJSON, projectRoot string) FuncCoverage {
	fc := FuncCoverage{
		File:      r.file,
		StartLine: r.startLine,
		EndLine:   r.endLine,
		Key:       fmt.Sprintf("%s:%d-%d", r.file, r.startLine, r.endLine),
	}

	if covData == nil {
		fc.CoveredPct = 100
		return fc
	}

	fileCov := findPyCoverageFile(covData, r.file, projectRoot)
	if fileCov == nil {
		fc.CoveredPct = 100
		return fc
	}

	countPyExecutedLines(fileCov, r, &fc)
	countPyMissingLines(fileCov, r, &fc)
	countPyExecutedBranches(fileCov, r, &fc)
	countPyMissingBranches(fileCov, r, &fc)

	if fc.TotalBlocks > 0 {
		fc.CoveredPct = float64(fc.CoveredBlocks) / float64(fc.TotalBlocks) * 100
	} else {
		fc.CoveredPct = 100
	}

	return fc
}
