//ff:type feature=runner type=model
//ff:what Defines the Runner interface for language-specific test execution
package runner

import "github.com/park-jun-woo/tsma/internal/match"

// Runner executes tests for a given language. The TestMatch carries the set of
// test files attributed to the source function plus, for content-aware (Go)
// matches, the explicit set of test functions to run. Non-Go implementations
// use m.Files[0] as the single test file and, when m.TestFuncs is nil, extract
// the functions from that file (preserving legacy behavior).
type Runner interface {
	Run(projectRoot string, m match.TestMatch) (*Result, error)
}
