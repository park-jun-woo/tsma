//ff:func feature=index type=helper control=sequence
//ff:what Builds a qualified name from package directory, receiver, and function name
package index

// buildQualifiedName constructs pkgDir.Receiver.Name or pkgDir.Name.
func buildQualifiedName(pkgDir, receiver, name string) string {
	if pkgDir == "" {
		if receiver != "" {
			return receiver + "." + name
		}
		return name
	}
	if receiver != "" {
		return pkgDir + "." + receiver + "." + name
	}
	return pkgDir + "." + name
}
