//ff:type feature=match type=model lang=go
//ff:what Points to one test function that references a source identifier
package match

// testRef points to one test function that references a source identifier.
type testRef struct {
	File     string // project-root-relative path of the _test.go file
	TestFunc string // name of the func TestXxx that references the identifier
}
