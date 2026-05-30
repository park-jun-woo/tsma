//ff:func feature=match type=helper control=sequence lang=go
//ff:what Adds a (name, receiver) distinguisher into the package source-receiver map
package match

// addNameReceiver records that the identifier name is declared with the given
// receiver distinguisher (a bare receiver type, or "" for a free function),
// creating the per-name distinguisher set on first use. Recording the same
// (name, receiver) pair twice is a no-op; recording a second distinct receiver
// for a name grows its set, which is how same-name-multiple is detected.
func addNameReceiver(r *PkgSourceReceivers, name, receiver string) {
	set := r.byName[name]
	if set == nil {
		set = make(map[string]struct{})
		r.byName[name] = set
	}
	set[receiver] = struct{}{}
}
