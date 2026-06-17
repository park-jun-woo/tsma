//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what pruneUnusedImports: file.Decls의 import 선언마다 filterImportSpecs로 미사용 spec을 솎고, 비게 된 import 블록은 통째로 제거한다. used 집합에 없는 패키지(단, blank `_`·dot `.`는 보존)가 미사용이다. decls→specs→if 3중 중첩을 피하려 spec 필터링은 filterImportSpecs로 추출(Q1 depth ≤2).

package tsmagate

import (
	"go/ast"
	"go/token"
)

// pruneUnusedImports filters each import declaration's specs through
// filterImportSpecs and drops any import block left empty. Spec filtering is
// extracted to keep nesting flat (decls→specs→if would be depth 3). The file's
// AST is mutated in place; gofmt re-emission happens in the caller.
func pruneUnusedImports(file *ast.File, used map[string]bool) {
	var decls []ast.Decl
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if ok && gd.Tok == token.IMPORT {
			gd.Specs = filterImportSpecs(gd.Specs, used)
		}
		if ok && gd.Tok == token.IMPORT && len(gd.Specs) == 0 {
			continue
		}
		decls = append(decls, decl)
	}
	file.Decls = decls
}
