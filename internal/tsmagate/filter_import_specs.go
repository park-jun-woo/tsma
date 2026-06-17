//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what filterImportSpecs: import spec 목록에서 보존할 것만 추린다. importIdent로 각 spec의 식별자를 구해 used에 있으면 보존하고, blank `_`·dot `.` import는 부작용/네임스페이스 주입 목적이라 사용 여부와 무관하게 항상 보존한다.

package tsmagate

import "go/ast"

// filterImportSpecs returns only the import specs to keep. importIdent yields each
// spec's package identifier; a spec is kept when that identifier is in used.
// Blank (`_`) and dot (`.`) imports are always kept regardless of use, since they
// exist for side effects or namespace injection rather than a qualifier.
func filterImportSpecs(specs []ast.Spec, used map[string]bool) []ast.Spec {
	var kept []ast.Spec
	for _, spec := range specs {
		id := importIdent(spec.(*ast.ImportSpec))
		if id == "_" || id == "." || used[id] {
			kept = append(kept, spec)
		}
	}
	return kept
}
