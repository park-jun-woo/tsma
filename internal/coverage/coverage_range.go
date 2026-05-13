//ff:type feature=coverage type=model
//ff:what Represents a source location range with start and end positions
package coverage

// coverageRange represents a source location range.
type coverageRange struct {
	Start coveragePosition `json:"start"`
	End   coveragePosition `json:"end"`
}
