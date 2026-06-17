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
	// Overlay carries the loop's non-invasive measurement plan: it maps absolute
	// virtual paths (a fresh _test.go inside the source package directory) to the
	// absolute backing file under .tsma/test that holds the generated test. When
	// non-empty, the Go runner/checker pass `go test -overlay <json> -vet=off` so
	// the test is compiled into the package without ever touching the source tree.
	// Empty for manual submit and non-Go matches (current disk-truth behavior).
	Overlay map[string]string
}
