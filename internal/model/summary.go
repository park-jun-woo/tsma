//ff:type feature=model type=model
//ff:what Holds aggregate counts of function statuses
package model

type Summary struct {
	Total int `json:"total"`
	Pass  int `json:"pass"`
	Done  int `json:"done"`
	Todo  int `json:"todo"`
}
