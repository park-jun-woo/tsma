//ff:type feature=match type=model lang=go
//ff:what Pairs a called identifier name with its statically-resolved receiver type
package match

// calledRef is a called identifier paired with the receiver type the call was
// made on, as resolved syntactically (go/ast only) from the call site. Receiver
// is the bare type name (pointer/value normalized) for a method call whose
// receiver could be statically determined, or "" (unknown) for a free-function
// call, an unresolvable method receiver, or any call that does not match the
// supported receiver-detection patterns.
type calledRef struct {
	Name     string
	Receiver string
}
