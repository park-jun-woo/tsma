//ff:func feature=index type=helper control=iteration dimension=1
//ff:what Returns true if the current item is inside a #[cfg(test)] scope or directly guarded by one
package index

// cfgTestActive reports whether a declaration should be treated as test-only,
// either because a #[cfg(test)] attribute directly precedes it (pending) or
// because an enclosing scope is marked cfgTest.
func cfgTestActive(scopes []rsScope, pending bool) bool {
	if pending {
		return true
	}
	for _, s := range scopes {
		if s.cfgTest {
			return true
		}
	}
	return false
}
