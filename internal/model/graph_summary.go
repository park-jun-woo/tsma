//ff:type feature=model type=model
//ff:what Holds aggregate counts of call graph properties
package model

type GraphSummary struct {
	Nodes       int `json:"nodes"`
	Edges       int `json:"edges"`
	EntryPoints int `json:"entry_points"`
	Dead        int `json:"dead"`
}
