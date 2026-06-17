//ff:func feature=smell type=helper control=sequence lang=rust
//ff:what rsPtrFieldAccess: reports whether a field_expression's `field` node is an `as_ptr` / `as_mut_ptr` raw-pointer extraction — the method a test calls to obtain a raw pointer it then dereferences. Node-based on the field_identifier, so a method merely named similarly in a string never matches. Narrow target (only the as_*_ptr family) so safe pointer-free code is untouched (false-positive zero, scoped to test code).
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsPtrAsPtr are the raw-pointer extraction methods.
var rsPtrAsPtr = map[string]bool{
	"as_ptr":     true,
	"as_mut_ptr": true,
}

// rsPtrFieldAccess returns the note + true when field is as_ptr/as_mut_ptr.
func rsPtrFieldAccess(field *treesitter.Node) (string, bool) {
	if field == nil || !rsPtrAsPtr[field.Text] {
		return "", false
	}
	return field.Text + "()", true
}
