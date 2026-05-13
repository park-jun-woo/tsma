//ff:type feature=chain type=model
//ff:what Stores parsed Go function metadata including AST body for chain tracing
package chain

import "go/ast"

// funcInfo stores parsed function metadata.
type funcInfo struct {
	name      string
	file      string
	startLine int
	endLine   int
	body      *ast.BlockStmt
}
