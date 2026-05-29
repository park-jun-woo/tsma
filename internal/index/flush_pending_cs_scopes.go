//ff:func feature=index type=helper control=iteration dimension=1 lang=csharp
//ff:what Pushes pending C# namespace/type names onto the scope stack at the given brace depth
package index

// flushPendingCsScopes moves every pending namespace/type name onto the scope
// stack, tagging each with the supplied brace depth (the depth in effect just
// before the opening brace is counted). It returns the updated scope stack and a
// cleared pending slice. This is invoked when an opening brace is seen so that
// declarations whose body brace sits on a later line (Allman style) are still
// attributed to the correct depth.
func flushPendingCsScopes(scopes []csScope, pending []string, depth int) ([]csScope, []string) {
	for _, name := range pending {
		scopes = append(scopes, csScope{depth: depth, typeName: name})
	}
	return scopes, nil
}
