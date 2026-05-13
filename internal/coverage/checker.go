//ff:type feature=coverage type=model
//ff:what Defines the Checker interface for language-specific branch coverage verification
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

// Checker verifies branch coverage for an endpoint.
type Checker interface {
	Check(projectRoot, testFile string, ep *model.Endpoint) (*Report, error)
}
