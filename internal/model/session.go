//ff:type feature=model type=model
//ff:what Represents the persistent state for a project
package model

import "time"

// Session represents the persistent state for a project.
type Session struct {
	Project   string     `json:"project"`
	Lang      string     `json:"lang"`
	Framework string     `json:"framework"`
	Created   time.Time  `json:"created"`
	Endpoints []Endpoint `json:"endpoints"`
	Summary   Summary    `json:"summary"`
}
