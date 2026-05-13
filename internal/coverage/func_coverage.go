//ff:type feature=coverage type=model
//ff:what Holds coverage data for a single function range
package coverage

// FuncCoverage holds coverage data for a single function range.
type FuncCoverage struct {
	Key            string // "file:startLine-endLine"
	File           string
	StartLine      int
	EndLine        int
	CoveredPct     float64
	TotalBlocks    int
	CoveredBlocks  int
	UncoveredLines []int
}
