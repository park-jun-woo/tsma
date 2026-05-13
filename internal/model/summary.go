//ff:type feature=model type=model
//ff:what Holds aggregate counts of endpoint statuses
package model

// Summary holds aggregate counts.
type Summary struct {
	Total   int `json:"total"`
	Done    int `json:"done"`
	Partial int `json:"partial"`
	Todo    int `json:"todo"`
}
