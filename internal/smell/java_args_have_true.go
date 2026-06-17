//ff:func feature=smell type=helper control=iteration dimension=1 lang=java
//ff:what javaArgsHaveTrue: reports whether an argument_list node has a direct `true` literal child — the discriminator that makes detectJavaSetAccessible fire on setAccessible(true) but not setAccessible(false). Looks only at direct children so a nested boolean inside a sub-expression argument does not accidentally trip it.
package smell

import "github.com/park-jun-woo/tsma/internal/treesitter"

// javaArgsHaveTrue reports whether args (an argument_list) contains a direct
// `true` literal argument.
func javaArgsHaveTrue(args *treesitter.Node) bool {
	if args == nil {
		return false
	}
	for _, c := range args.Children {
		if c.Type == "true" {
			return true
		}
	}
	return false
}
