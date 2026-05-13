//ff:type feature=model type=model
//ff:what Represents a single function to be tested
package model

type Function struct {
	QualifiedName     string  `json:"qualified_name"`
	Name              string  `json:"name"`
	File              string  `json:"file"`
	StartLine         int     `json:"start_line"`
	EndLine           int     `json:"end_line"`
	Exported          bool    `json:"exported"`
	Status            string  `json:"status"`
	TestFile          string  `json:"test_file,omitempty"`
	CoveragePct       float64 `json:"coverage_pct,omitempty"`
	UncoveredBranches []int   `json:"uncovered_branches,omitempty"`
	RetryCount        int     `json:"retry_count,omitempty"`
}
