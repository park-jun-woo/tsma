//ff:type feature=coverage type=model
//ff:what Represents a branch entry in an istanbul coverage report
package coverage

// coverageBranch represents a branch entry in the coverage report.
type coverageBranch struct {
	Loc       coverageRange   `json:"loc"`
	Type      string          `json:"type"`
	Locations []coverageRange `json:"locations"`
}
