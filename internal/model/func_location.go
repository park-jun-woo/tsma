//ff:type feature=model type=model
//ff:what Represents a function's location in source code
package model

// FuncLocation represents a function's location in source code.
type FuncLocation struct {
	File      string `json:"file"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}
