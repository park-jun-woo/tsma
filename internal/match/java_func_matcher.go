//ff:type feature=match type=implementation lang=java
//ff:what Content-aware FuncMatcher that attributes Java tests by analyzing which test files (in the JUnit src/test mirror) call a function (the M-001 analogue of GoFuncMatcher/TypeScriptFuncMatcher). Filename matching via JavaMatcher (FooTest.java/FooTests.java) is retained as the last-resort fallback, never overriding a content match.
package match

// JavaFuncMatcher attributes tests to a Java method/constructor by analyzing
// which test files reference it (content-aware via tree-sitter), rather than by
// the FooTest.java file-name convention alone. Same shape as GoFuncMatcher.
type JavaFuncMatcher struct{}
