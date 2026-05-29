//ff:type feature=match type=model lang=go
//ff:what Maps source identifier names to the Go test functions that reference them
package match

// PkgTestIndex is a content-aware index for a single Go package directory.
// Keys are bare identifier names (function names / method names) referenced by
// the package's test functions. Each key maps to the test functions that call
// or reference that identifier (directly, or via a same-file 1-hop helper).
type PkgTestIndex struct {
	refs map[string][]testRef
}
