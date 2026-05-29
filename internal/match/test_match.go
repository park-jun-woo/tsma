//ff:type feature=match type=model
//ff:what Holds the test files and test functions attributed to one source function
package match

// TestMatch is the result of attributing tests to a single source function.
// Files holds project-root-relative paths of the _test.go files that cover the
// function; TestFuncs holds the names of the test functions to execute. A nil
// or empty TestFuncs means "run every test in Files" (left for the runner to
// resolve), which is how non-Go fallback matches report their result.
type TestMatch struct {
	Files     []string
	TestFuncs []string
}
