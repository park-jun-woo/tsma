//ff:type feature=coverage type=model
//ff:what Represents coverage data for a single Python file
package coverage

// pyCoverageFile represents coverage data for a single file.
type pyCoverageFile struct {
	ExecutedLines    []int   `json:"executed_lines"`
	MissingLines     []int   `json:"missing_lines"`
	ExecutedBranches [][]int `json:"executed_branches"`
	MissingBranches  [][]int `json:"missing_branches"`
}
