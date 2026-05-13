//ff:type feature=endpoint type=factory
//ff:what Interface for extracting endpoints from a project
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

// Detector extracts endpoints from a project.
type Detector interface {
	Detect(projectRoot string) ([]model.Endpoint, error)
}
