//ff:type feature=match type=implementation lang=rust
//ff:what Content-aware FuncMatcher that attributes Rust tests by analyzing which test files (the in-file #[cfg(test)] module or a tests/*.rs integration test) call a function — the M-001 analogue of GoFuncMatcher/CsFuncMatcher. Filename matching via RsMatcher (source file itself for in-file tests, else tests/<name>.rs) is retained as the last-resort fallback, never overriding a content match. Covers non-pub functions (the in-file module sees them), which the tests/ integration path alone cannot.
package match

// RsFuncMatcher attributes tests to a Rust function by analyzing which test
// files reference it (content-aware via tree-sitter), rather than by the in-file
// / tests/ filename convention alone. Same shape as CsFuncMatcher.
type RsFuncMatcher struct{}
