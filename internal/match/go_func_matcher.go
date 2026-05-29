//ff:type feature=match type=implementation lang=go
//ff:what Content-aware FuncMatcher that attributes Go tests by analyzing call sites
package match

// GoFuncMatcher attributes tests to a Go function by analyzing which test
// functions reference it (content-aware), rather than by file-name convention.
type GoFuncMatcher struct{}
