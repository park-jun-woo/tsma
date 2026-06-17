//ff:func feature=smell type=helper control=sequence lang=rust
//ff:what rsFnIsUnsafe: reports whether a function_item is declared `unsafe` by inspecting its `function_modifiers` child leaf text. Node-based — the modifiers node carries only the keyword tokens (unsafe/async/const), so an "unsafe" substring inside an identifier or string never matches. Used by detectRsUnsafe to flag an `unsafe fn` test helper alongside bare unsafe blocks.
package smell

import (
	"strings"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// rsFnIsUnsafe returns true when a function_item carries the unsafe modifier.
func rsFnIsUnsafe(fn *treesitter.Node) bool {
	mods := fn.ChildByType("function_modifiers")
	if mods == nil {
		return false
	}
	return strings.Contains(mods.Text, "unsafe")
}
