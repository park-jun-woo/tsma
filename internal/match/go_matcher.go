//ff:type feature=match type=implementation
//ff:what Go function-test matcher that scans *_test.go files via AST
package match

// GoMatcher finds test files for Go functions.
type GoMatcher struct{}
