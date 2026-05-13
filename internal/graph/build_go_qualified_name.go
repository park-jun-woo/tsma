//ff:func feature=graph type=helper control=sequence
//ff:what Constructs the qualified name for a Go ast.FuncDecl with optional receiver
package graph

import "go/ast"

// buildGoQualifiedName constructs the qualified name for an ast.FuncDecl.
func buildGoQualifiedName(fd *ast.FuncDecl, pkgDir string) string {
	name := fd.Name.Name
	var receiver string
	if fd.Recv != nil && len(fd.Recv.List) > 0 {
		receiver = extractRecvType(fd.Recv.List[0].Type)
	}

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
