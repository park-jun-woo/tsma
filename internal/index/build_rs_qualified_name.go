//ff:func feature=index type=helper control=iteration dimension=1
//ff:what Builds a Rust qualified name from package dir, enclosing scope stack, and function name
package index

import "strings"

// buildRsQualifiedName constructs "pkgDir.mod::Recv::name" style qualified names.
// Module and impl-receiver segments from the scope stack are joined with "::",
// then combined with the package directory via the shared buildQualifiedName.
func buildRsQualifiedName(pkgDir string, scopes []rsScope, name string) string {
	var segs []string
	for _, s := range scopes {
		if s.module != "" {
			segs = append(segs, s.module)
		}
		if s.receiver != "" {
			segs = append(segs, s.receiver)
		}
	}
	receiver := strings.Join(segs, "::")
	return buildQualifiedName(pkgDir, receiver, name)
}
