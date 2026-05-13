//ff:func feature=chain type=helper control=sequence
//ff:what Builds the index key for a Go function declaration with optional receiver prefix
package chain

import "go/ast"

// buildFuncKey builds the index key for a function declaration.
func buildFuncKey(fd *ast.FuncDecl, pkgDir string) string {
	if fd.Recv != nil && len(fd.Recv.List) > 0 {
		recvType := extractRecvType(fd.Recv.List[0].Type)
		return pkgDir + "." + recvType + "." + fd.Name.Name
	}
	return pkgDir + "." + fd.Name.Name
}
