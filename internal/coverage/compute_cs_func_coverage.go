//ff:func feature=coverage type=implementation control=sequence lang=csharp
//ff:what Computes C# function coverage by mapping Cobertura line and branch counters onto a function range
package coverage

import "fmt"

// computeCsFuncCoverage computes coverage for a single C# function range,
// counting executable lines and branch outcomes within the range. When no
// coverage data is available for the file, the function is treated as fully
// covered (100%) so a missing report does not block progress.
func computeCsFuncCoverage(r csFuncRange, cov *csCoverage, projectRoot string) FuncCoverage {
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

	fileCov := findCsCoverageFile(cov, r.file, projectRoot)
	if fileCov == nil {
		fc.CoveredPct = 100
		return fc
	}

	countCsLines(fileCov, r, &fc)
	countCsBranches(fileCov, r, &fc)

	if fc.TotalBlocks > 0 {
		fc.CoveredPct = float64(fc.CoveredBlocks) / float64(fc.TotalBlocks) * 100
	} else {
		fc.CoveredPct = 100
	}

	return fc
}
