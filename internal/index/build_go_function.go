//ff:func feature=index type=helper control=sequence
//ff:what Constructs a model.Function from an ast.FuncDecl with qualified name
package index

import (
	"go/ast"
	"go/token"
	"unicode"

	"github.com/park-jun-woo/tsma/internal/model"
)

// buildGoFunction constructs a model.Function from an ast.FuncDecl.
func buildGoFunction(fd *ast.FuncDecl, fset *token.FileSet, relPath, pkgDir string) model.Function {
	name := fd.Name.Name

	var receiver string
	if fd.Recv != nil && len(fd.Recv.List) > 0 {
		receiver = extractReceiver(fd.Recv.List[0].Type)
	}

	qualifiedName := buildQualifiedName(pkgDir, receiver, name)
	exported := len(name) > 0 && unicode.IsUpper(rune(name[0]))

	return model.Function{
		QualifiedName: qualifiedName,
		Name:          name,
		Receiver:      receiver,
		File:          relPath,
		StartLine:     fset.Position(fd.Pos()).Line,
		EndLine:       fset.Position(fd.End()).Line,
		Exported:      exported,
		Status:        model.StatusTodo,
	}
}
