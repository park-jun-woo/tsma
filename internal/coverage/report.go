//ff:type feature=coverage type=model
//ff:what Holds coverage results for an endpoint's call chain
package coverage

// Report holds coverage results for an endpoint's chain.
type Report struct {
	Funcs      []FuncCoverage
	AllCovered bool
	TotalPct   float64
	Uncovered  []UncoveredBranch
}
