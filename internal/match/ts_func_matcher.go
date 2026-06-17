//ff:type feature=match type=implementation lang=typescript
//ff:what Content-aware FuncMatcher that attributes TypeScript tests by analyzing which test files call a function (the M-001 analogue of GoFuncMatcher). Filename matching via TSMatcher is retained as the last-resort fallback, never overriding a content match.
package match

// TypeScriptFuncMatcher attributes tests to a TS/JS function by analyzing which
// test files reference it (content-aware via tree-sitter), rather than by
// file-name convention alone. Same shape as GoFuncMatcher.
type TypeScriptFuncMatcher struct{}
