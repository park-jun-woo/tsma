//ff:type feature=coverage type=model
//ff:what Represents a function entry in an istanbul coverage report
package coverage

// coverageFunction represents a function entry in the coverage report.
type coverageFunction struct {
	Name string        `json:"name"`
	Loc  coverageRange `json:"loc"`
	Decl coverageRange `json:"decl"`
}
