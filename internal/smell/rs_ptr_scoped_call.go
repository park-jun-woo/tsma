//ff:func feature=smell type=helper control=sequence lang=rust
//ff:what rsPtrScopedCall: reports whether a call_expression's `function` node is a `std::ptr::read` / `core::ptr::write` style raw-pointer access — a scoped_identifier whose rightmost `name` is read/write/read_unaligned/write_unaligned and whose `path` segment's name is exactly `ptr`. Requiring the `ptr` qualifier keeps a plain `reader.read()` / `buf.write()` (a field_expression, no `ptr` path) from ever matching (false-positive zero). Returns the call note on a match.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// rsPtrReadWrite are the std::ptr free functions that force a raw memory access.
var rsPtrReadWrite = map[string]bool{
	"read": true, "write": true,
	"read_unaligned": true, "write_unaligned": true,
	"read_volatile": true, "write_volatile": true,
}

// rsPtrScopedCall returns the note + true when fn is a std::ptr::read/write call.
func rsPtrScopedCall(fn *treesitter.Node) (string, bool) {
	if fn == nil || fn.Type != "scoped_identifier" {
		return "", false
	}
	name := fn.ChildByField("name")
	if name == nil || !rsPtrReadWrite[name.Text] {
		return "", false
	}
	path := fn.ChildByField("path")
	if path == nil || path.Type != "scoped_identifier" {
		return "", false
	}
	seg := path.ChildByField("name")
	if seg == nil || seg.Text != "ptr" {
		return "", false
	}
	return "ptr::" + name.Text + "()", true
}
