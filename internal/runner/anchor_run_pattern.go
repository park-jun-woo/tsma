//ff:func feature=runner type=util control=sequence lang=go
//ff:what Builds an exact-match anchored -run regex from a set of test function names
package runner

import "strings"

// AnchorRunPattern builds an exact-match `go test -run` regex from test function
// names: `^(TestA|TestB)$`. Anchoring prevents prefix over-execution (so
// `TestGenerate` does not also match `TestGenerateBytes`). It returns "" for an
// empty set so callers omit -run and run the whole package.
func AnchorRunPattern(testFuncs []string) string {
	if len(testFuncs) == 0 {
		return ""
	}
	if len(testFuncs) == 1 {
		return "^" + testFuncs[0] + "$"
	}
	return "^(" + strings.Join(testFuncs, "|") + ")$"
}
