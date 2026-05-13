//ff:type feature=match type=factory
//ff:what Defines the Matcher interface for function-test file matching
package match

import "github.com/park-jun-woo/tsma/internal/model"

// Matcher finds the test file that covers a given function.
type Matcher interface {
	Match(projectRoot string, fn *model.Function) (testFile string, found bool)
}
