//ff:type feature=model type=model
//ff:what Holds aggregate counts of function statuses
package model

type Summary struct {
	Total   int `json:"total"`
	Done    int `json:"done"`
	Partial int `json:"partial"`
	Todo    int `json:"todo"`
}
