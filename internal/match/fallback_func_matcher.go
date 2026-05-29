//ff:type feature=match type=implementation
//ff:what Adapts a file-name based Matcher to the FuncMatcher interface for non-Go languages
package match

// fallbackFuncMatcher adapts an existing file-name based Matcher to the
// FuncMatcher interface, preserving the legacy single-file matching behavior for
// languages that have no content-aware implementation. It returns the single
// matched test file with a nil TestFuncs (meaning "run every test in the file");
// it deliberately does not extract test function names so that the match package
// never imports the runner package (avoiding an import cycle). The runner fills
// in TestFuncs from the file when TestFuncs is nil.
type fallbackFuncMatcher struct {
	inner Matcher
}
