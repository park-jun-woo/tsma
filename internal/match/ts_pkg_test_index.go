//ff:type feature=match type=model lang=typescript
//ff:what Maps bare source identifier names to the TS/JS test files that call them — the content-aware index for one package directory. Coarser than the Go index (file granularity, no receiver), because jest/vitest run whole files; TestFuncs is therefore left nil downstream.
package match

// TSPkgTestIndex is a content-aware index for a single TS/JS package directory.
// Keys are bare identifier names (function / method / constructor names) called
// by the directory's test files; each maps to the test files that reference it.
type TSPkgTestIndex struct {
	refs map[string][]string
}
