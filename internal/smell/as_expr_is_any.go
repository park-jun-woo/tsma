//ff:func feature=smell type=helper control=sequence lang=typescript
//ff:what asExprIsAny: reports whether an as_expression casts to the predefined `any` type (not a named type, which is a legitimate narrowing cast) — the discriminator that keeps `as T` from firing TS-REFL-TS-001.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// asExprIsAny reports whether an as_expression casts to the predefined `any`
// type.
func asExprIsAny(ae *treesitter.Node) bool {
	t := ae.ChildByType("predefined_type")
	return t != nil && t.Text == "any"
}
