//ff:func feature=gate type=helper control=sequence
//ff:what importIdent: import spec이 파일 안에서 한정자로 쓰일 식별자를 돌려준다. 명시 별칭(blank `_`·dot `.` 포함)이 있으면 그 이름을, 없으면 import 경로의 마지막 요소(base)를 쓴다. LLM 테스트의 import는 전부 stdlib라 경로 base가 패키지명과 정확히 일치한다.

package tsmagate

import (
	"go/ast"
	"path"
	"strconv"
)

// importIdent returns the identifier under which an import spec would be
// qualified in the file. An explicit name (including blank `_` and dot `.`) wins;
// otherwise the last element of the import path is used. For LLM-generated tests
// every import is stdlib, so the path base equals the package name exactly.
func importIdent(spec *ast.ImportSpec) string {
	if spec.Name != nil {
		return spec.Name.Name
	}
	p, err := strconv.Unquote(spec.Path.Value)
	if err != nil {
		return spec.Path.Value
	}
	return path.Base(p)
}
