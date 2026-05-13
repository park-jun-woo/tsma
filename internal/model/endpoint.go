//ff:type feature=model type=model
//ff:what Represents a single API endpoint to be tested
package model

// Endpoint represents a single API endpoint to be tested.
type Endpoint struct {
	Name              string             `json:"name"`
	Method            string             `json:"method"`
	Path              string             `json:"path"`
	Handler           FuncLocation       `json:"handler"`
	Chain             []ChainEntry       `json:"chain"`
	Status            string             `json:"status"`
	TestFile          string             `json:"test_file,omitempty"`
	Coverage          map[string]float64 `json:"coverage,omitempty"`
	UncoveredBranches []int              `json:"uncovered_branches,omitempty"`
}
