//ff:type feature=model type=model
//ff:what Represents the persistent state for a project
package model

import "time"

type Session struct {
	Project   string     `json:"project"`
	Lang      string     `json:"lang"`
	Created   time.Time  `json:"created"`
	Functions []Function `json:"functions"`
	Summary   Summary    `json:"summary"`
}
