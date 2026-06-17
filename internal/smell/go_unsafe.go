//ff:func feature=smell type=helper control=iteration dimension=1 level=error
//ff:what detectUnsafe: TS-REFL-001. go/ast로 (a) `import "unsafe"` ImportSpec와 (b) `unsafe.` SelectorExpr(X가 Ident "unsafe")를 정확히 탐지한다. unsafe는 메모리로 비공개 상태를 강제 조작하는 강한 cheese라 Review로 표면화한다. 문자열 리터럴 "unsafe"나 주석은 AST 노드가 아니므로 발화하지 않는다(위양성 0).

package smell

import (
	"go/ast"
	"go/token"
)

// detectUnsafe finds TS-REFL-001 violations: importing or using the unsafe
// package in a test. It matches the ImportSpec whose path is "unsafe" and any
// selector whose base identifier is unsafe (e.g. unsafe.Pointer). Matching is on
// AST nodes, so an "unsafe" string literal or a comment mentioning unsafe never
// fires.
func detectUnsafe(file *ast.File, fset *token.FileSet, path string) []Finding {
	var findings []Finding

	for _, imp := range file.Imports {
		if imp.Path != nil && imp.Path.Value == `"unsafe"` {
			findings = append(findings, Finding{
				Rule: "TS-REFL-001",
				File: path,
				Line: fset.Position(imp.Pos()).Line,
				Note: `import "unsafe"`,
			})
		}
	}

	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "unsafe" {
			findings = append(findings, Finding{
				Rule: "TS-REFL-001",
				File: path,
				Line: fset.Position(sel.Pos()).Line,
				Note: "unsafe." + sel.Sel.Name,
			})
		}
		return true
	})

	return findings
}
