//ff:type feature=match type=model lang=csharp
//ff:what Maps bare source identifier names (method/constructor/property names) to the test files (in the *.Tests project) that call them — the content-aware index for one C# source directory's test candidates. Coarser than the Go index (file granularity, no receiver), because dotnet runs whole test classes; TestFuncs is therefore left nil downstream.
package match

// CsPkgTestIndex is a content-aware index for a single C# source directory's
// test candidates. Keys are bare identifier names (method / constructor /
// property names) called by the candidate *Tests.cs files; each maps to the test
// files that reference it.
type CsPkgTestIndex struct {
	refs map[string][]string
}
