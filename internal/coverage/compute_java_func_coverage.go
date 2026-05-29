//ff:func feature=coverage type=implementation control=sequence
//ff:what Computes Java function coverage by mapping JaCoCo line and branch counters onto a function range
package coverage

import "fmt"

// computeJavaFuncCoverage computes coverage for a single Java function range,
// counting executable lines and branch outcomes within the range.
func computeJavaFuncCoverage(r javaFuncRange, cov *jacocoCoverage, projectRoot string) FuncCoverage {
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

	fileCov := findJavaCoverageFile(cov, r.file, projectRoot)
	if fileCov == nil {
		fc.CoveredPct = 100
		return fc
	}

	countJavaLines(fileCov, r, &fc)
	countJavaBranches(fileCov, r, &fc)

	if fc.TotalBlocks > 0 {
		fc.CoveredPct = float64(fc.CoveredBlocks) / float64(fc.TotalBlocks) * 100
	} else {
		fc.CoveredPct = 100
	}

	return fc
}
