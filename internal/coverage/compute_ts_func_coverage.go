//ff:func feature=coverage type=implementation control=sequence
//ff:what Matches TS function range against istanbul coverage data, computes covered/uncovered stats
package coverage

import "fmt"

// computeTSFuncCoverage computes statement and branch coverage for a specific function range.
func computeTSFuncCoverage(r tsFuncRange, coverageData map[string]coverageFinalEntry, projectRoot string) FuncCoverage {
	fc := FuncCoverage{
		File:      r.file,
		StartLine: r.startLine,
		EndLine:   r.endLine,
		Key:       fmt.Sprintf("%s:%d-%d", r.file, r.startLine, r.endLine),
	}

	entry := findCoverageEntry(r.file, coverageData, projectRoot)
	if entry == nil {
		fc.CoveredPct = 100
		return fc
	}

	countTSStatements(entry, r, &fc)
	countTSBranches(entry, r, &fc)

	if fc.TotalBlocks > 0 {
		fc.CoveredPct = float64(fc.CoveredBlocks) / float64(fc.TotalBlocks) * 100
	} else {
		fc.CoveredPct = 100
	}

	return fc
}
