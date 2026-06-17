//ff:type feature=match type=implementation lang=csharp
//ff:what Content-aware FuncMatcher that attributes C# tests by analyzing which test files (in the parallel *.Tests project) call a function (the M-001 analogue of GoFuncMatcher/JavaFuncMatcher). Filename matching via CsMatcher (FooTests.cs/FooTest.cs) is retained as the last-resort fallback, never overriding a content match.
package match

// CsFuncMatcher attributes tests to a C# method/constructor/property by
// analyzing which test files reference it (content-aware via tree-sitter),
// rather than by the FooTests.cs file-name convention alone. Same shape as
// JavaFuncMatcher.
type CsFuncMatcher struct{}
