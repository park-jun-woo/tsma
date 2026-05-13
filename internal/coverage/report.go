//ff:type feature=coverage type=model
//ff:what Holds coverage results for a function
package coverage

// Report holds coverage results for a function.
type Report struct {
	Funcs      []FuncCoverage
	AllCovered bool
	TotalPct   float64
	Uncovered  []UncoveredBranch
}
