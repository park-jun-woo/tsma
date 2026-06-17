//ff:type feature=match type=implementation lang=python
//ff:what Content-aware FuncMatcher that attributes Python tests by analyzing which test files reference a function (the M-001 analogue of GoFuncMatcher, via the built-in ast). Filename matching via PyMatcher is retained as the last-resort fallback, never overriding a content match.
package match

// PythonFuncMatcher attributes tests to a Python function by analyzing which
// test files reference it (content-aware via the ast subprocess), rather than by
// file-name convention alone. Same shape as TypeScriptFuncMatcher.
type PythonFuncMatcher struct{}
