//ff:type feature=coverage type=model
//ff:what Defines the Checker interface for language-specific branch coverage verification
package coverage

import (
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// Checker verifies branch coverage for a function. The TestMatch carries the
// set of test files attributed to the function plus, for content-aware (Go)
// matches, the explicit set of test functions to run. Non-Go implementations
// use m.Files[0] as the single test file (behavior unchanged).
type Checker interface {
	Check(projectRoot string, m match.TestMatch, fn *model.Function) (*Report, error)
}
