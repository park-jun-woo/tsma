//ff:func feature=match type=helper control=selection lang=go
//ff:what Decides whether one test ref belongs to a method given its receiver
package match

// keepRefForReceiver decides whether a single test ref r should be attributed to
// a method whose receiver type is fnReceiver. multiple reports whether the
// method's name is same-name-multiple in the package. The rule:
//
//   - an exact receiver match is always kept;
//   - an unknown-receiver ref ("") is kept only when the name is
//     same-name-single (no other declaration shares it), preserving the
//     pre-existing correct single-receiver behavior; when same-name-multiple it
//     is dropped to avoid mis-attributing another type's test;
//   - any other (concrete but mismatched) receiver is dropped.
func keepRefForReceiver(r testRef, fnReceiver string, multiple bool) bool {
	switch {
	case r.Receiver == fnReceiver:
		return true
	case r.Receiver == "":
		return !multiple
	default:
		return false
	}
}
