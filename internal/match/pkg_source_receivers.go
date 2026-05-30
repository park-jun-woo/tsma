//ff:type feature=match type=model lang=go
//ff:what Maps a package's declared names to the set of receiver types they appear on
package match

// PkgSourceReceivers records, for a single Go package directory, the set of
// distinguishers each declared identifier name appears with across the
// package's non-test .go sources. For a method the distinguisher is its bare
// receiver type name; for a free function it is "" (the empty receiver). A name
// whose distinguisher set has more than one element is "same-name-multiple":
// either the same method name declared on several receiver types, or a free
// function and a method sharing a name. This drives the matching policy in
// MatchFuncByName — an unknown-receiver test reference is attributed only when
// its name is same-name-single (set size 1), and dropped when same-name-
// multiple (set size > 1), to avoid mis-attribution. It is built once per
// package directory (BuildPkgSourceReceivers), in parallel with the test index,
// and is a separate structure from PkgTestIndex (source declarations, not test
// refs).
type PkgSourceReceivers struct {
	byName map[string]map[string]struct{}
}
