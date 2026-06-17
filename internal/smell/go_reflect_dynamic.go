//ff:func feature=smell type=helper control=sequence level=error
//ff:what detectReflectDynamic: TS-REFL-002. *ast.SelectorExpr의 Sel.Name ∈ {MethodByName, FieldByName}만 표적으로 비공개 메서드/필드 동적 침투를 탐지한다. reflect.DeepEqual/TypeOf/ValueOf는 이 집합에 없으므로 절대 발화하지 않는다(가장 흔한 위양성원 차단). 1차는 메서드명 매칭으로 충분(두 이름은 사실상 reflect 전용).

package smell

import (
	"go/ast"
	"go/token"
)

// dynamicReflectSel is the set of reflect selector names that reach into private
// methods/fields dynamically. DeepEqual/TypeOf/ValueOf are deliberately absent:
// they are legitimate idioms and must never fire (rulebook §6 false-positive
// zero).
var dynamicReflectSel = map[string]bool{
	"MethodByName": true,
	"FieldByName":  true,
}

// detectReflectDynamic finds TS-REFL-002 violations: a selector ending in
// MethodByName or FieldByName (e.g. v.MethodByName("x")). It matches on the
// selector's Sel.Name only — 1st-pass method-name matching, since these two
// names are effectively reflect-exclusive. It never matches DeepEqual/TypeOf/
// ValueOf.
func detectReflectDynamic(file *ast.File, fset *token.FileSet, path string) []Finding {
	var findings []Finding

	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if dynamicReflectSel[sel.Sel.Name] {
			findings = append(findings, Finding{
				Rule: "TS-REFL-002",
				File: path,
				Line: fset.Position(sel.Sel.Pos()).Line,
				Note: "reflect dynamic ." + sel.Sel.Name,
			})
		}
		return true
	})

	return findings
}
