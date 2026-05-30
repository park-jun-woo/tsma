//ff:func feature=match type=helper control=sequence lang=go
//ff:what Marks a local variable name ambiguous so its type resolves to unknown
package match

// poison marks a local variable name as ambiguous: it removes any tracked type
// for the name and records it as poisoned so a later assignment cannot restore
// a type for it. Used by localVarTypes to enforce the single-binding rule
// (reassignment, branch/loop assignment, or multi-assign all drop the name).
func poison(types map[string]string, poisoned map[string]struct{}, name string) {
	delete(types, name)
	poisoned[name] = struct{}{}
}
