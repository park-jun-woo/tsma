//ff:func feature=index type=helper control=iteration dimension=1 lang=csharp
//ff:what Pops C# namespace/type scopes whose body has closed at the current brace depth
package index

// popClosedCsScopes removes scopes from the top of the stack whose body has
// closed, i.e. the current brace depth has fallen back to (or below) the depth
// at which the scope's body was opened.
func popClosedCsScopes(scopes []csScope, depth int) []csScope {
	for len(scopes) > 0 && depth <= scopes[len(scopes)-1].depth {
		scopes = scopes[:len(scopes)-1]
	}
	return scopes
}
