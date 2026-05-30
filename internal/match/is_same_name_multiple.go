//ff:func feature=match type=helper control=selection lang=go
//ff:what Reports whether a name is declared on more than one receiver in the package
package match

// isSameNameMultiple reports whether name is declared with more than one
// distinguisher in the package — i.e. the same method name on several receiver
// types, or a free function and a method sharing the name. A name with exactly
// one distinguisher (or one not present in the source map, e.g. when the
// package sources could not be parsed) is treated as same-name-single, so
// matching stays conservative toward the existing single-receiver behavior and
// never drops a reference solely because the source map is incomplete. A nil
// receiver record reports false (single) for every name.
func (r *PkgSourceReceivers) isSameNameMultiple(name string) bool {
	if r == nil || r.byName == nil {
		return false
	}
	return len(r.byName[name]) > 1
}
