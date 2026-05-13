//ff:type feature=model type=model
//ff:what Represents a single function in the call chain
package model

// ChainEntry represents a single function in the call chain.
type ChainEntry struct {
	Func      string `json:"func"`
	File      string `json:"file,omitempty"`
	StartLine int    `json:"start_line,omitempty"`
	EndLine   int    `json:"end_line,omitempty"`
	Boundary  string `json:"boundary,omitempty"`
}
