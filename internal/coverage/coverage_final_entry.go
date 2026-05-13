//ff:type feature=coverage type=model
//ff:what Represents a single file entry in coverage-final.json istanbul format
package coverage

// coverageFinalEntry represents a single file entry in coverage-final.json (istanbul format).
type coverageFinalEntry struct {
	StatementMap map[string]coverageRange    `json:"statementMap"`
	S            map[string]int              `json:"s"`
	BranchMap    map[string]coverageBranch   `json:"branchMap"`
	B            map[string][]int            `json:"b"`
	FnMap        map[string]coverageFunction `json:"fnMap"`
	F            map[string]int              `json:"f"`
}
