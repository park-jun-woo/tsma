//ff:func feature=match type=helper control=sequence lang=go
//ff:what Reports whether a function name is a Go test entrypoint (TestXxx)
package match

// isTestFuncName reports whether name is a Go test entrypoint of the form
// TestXxx. Bare "Test" is excluded (go test requires a following character).
func isTestFuncName(name string) bool {
	if len(name) <= len("Test") {
		return false
	}
	return name[:len("Test")] == "Test"
}
