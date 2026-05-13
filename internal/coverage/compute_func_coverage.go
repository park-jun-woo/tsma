//ff:func feature=coverage type=implementation control=iteration dimension=1
//ff:what Computes Go function coverage by matching cover blocks against a function range
package coverage

import "fmt"

func computeFuncCoverage(r funcRange, blocks []coverBlock, projectRoot string) FuncCoverage {
	fc := FuncCoverage{
		File:      r.file,
		StartLine: r.startLine,
		EndLine:   r.endLine,
		Key:       fmt.Sprintf("%s:%d-%d", r.file, r.startLine, r.endLine),
	}

	for _, b := range blocks {
		if !overlaps(b.file, r.file, b.startLine, b.endLine, r.startLine, r.endLine) {
			continue
		}
		fc.TotalBlocks++
		if b.count > 0 {
			fc.CoveredBlocks++
		} else {
			fc.UncoveredLines = append(fc.UncoveredLines, b.startLine)
		}
	}

	if fc.TotalBlocks > 0 {
		fc.CoveredPct = float64(fc.CoveredBlocks) / float64(fc.TotalBlocks) * 100
	} else {
		fc.CoveredPct = 100
	}

	return fc
}
