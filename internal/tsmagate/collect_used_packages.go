//ff:func feature=gate type=helper control=sequence
//ff:what collectUsedPackages: AST를 순회해 패키지 한정자로 쓰인 식별자 집합을 모은다. `pkg.Sym`은 *ast.SelectorExpr{X: Ident{pkg}}로 파싱되므로 X가 단순 식별자인 SelectorExpr마다 그 이름을 used에 넣는다. tidyGoSource가 import 미사용 판정에 쓴다.

package tsmagate

import "go/ast"

// collectUsedPackages walks the file and returns the set of identifiers used as a
// package qualifier. A reference `pkg.Sym` parses to *ast.SelectorExpr whose X is
// the *ast.Ident "pkg", so every selector with a bare-identifier X contributes
// its name. The result drives the unused-import decision in tidyGoSource.
func collectUsedPackages(file *ast.File) map[string]bool {
	used := make(map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				used[ident.Name] = true
			}
		}
		return true
	})
	return used
}
