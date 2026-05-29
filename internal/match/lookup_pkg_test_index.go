//ff:func feature=match type=helper control=sequence lang=go
//ff:what Looks up the test references for a bare identifier name in the index
package match

// Lookup returns the test references for a bare identifier name (e.g.
// model.Function.Name). Returns nil and false when the identifier is not
// referenced by any test in the package.
func (idx *PkgTestIndex) Lookup(name string) ([]testRef, bool) {
	if idx == nil || idx.refs == nil {
		return nil, false
	}
	refs, ok := idx.refs[name]
	if !ok || len(refs) == 0 {
		return nil, false
	}
	return refs, true
}
