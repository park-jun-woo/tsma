//ff:type feature=model type=model
//ff:what Represents a call graph edge with ambiguity flag
package model

type Edge struct {
	Target    string `json:"target"`
	Ambiguous bool   `json:"ambiguous"`
}
