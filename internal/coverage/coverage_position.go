//ff:type feature=coverage type=model
//ff:what Represents a line/column position in source code
package coverage

// coveragePosition represents a line/column position.
type coveragePosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}
