//ff:type feature=model type=model
//ff:what Represents the persistent session state for a project
package model

import "time"

type Session struct {
	Project      string     `json:"project"`
	Lang         string     `json:"lang"`
	CheckedAt    time.Time  `json:"checked_at"`
	CurrentIndex int        `json:"current_index"`
	Functions    []Function `json:"functions"`
	Summary      Summary    `json:"summary"`
}
