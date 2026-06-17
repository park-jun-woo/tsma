//ff:type feature=match type=model lang=java
//ff:what Maps bare source identifier names (method/constructor names) to the JUnit test files that call them — the content-aware index for one Java package's test mirror directory. Coarser than the Go index (file granularity, no receiver), because JUnit runs whole classes; TestFuncs is therefore left nil downstream.
package match

// JavaPkgTestIndex is a content-aware index for a single Java package's test
// mirror directory. Keys are bare identifier names (method / constructor names)
// called by the directory's *Test.java files; each maps to the test files that
// reference it.
type JavaPkgTestIndex struct {
	refs map[string][]string
}
