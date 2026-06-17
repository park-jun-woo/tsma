//ff:func feature=smell type=helper control=iteration dimension=1 level=error
//ff:what linknameInGroup: detectLinkname의 내부 헬퍼. 한 CommentGroup의 주석들을 돌며 c.Text가 "//go:linkname"로 시작하는 지시자만 Finding(TS-REFL-003)으로 낸다. detectLinkname에서 분리해 주석 순회 중첩을 평탄화한다(Q1 depth ≤2). prefix가 다른 일반 주석은 발화하지 않는다(위양성 0).

package smell

import (
	"go/ast"
	"go/token"
	"strings"
)

// linknameInGroup scans one comment group for a //go:linkname directive and
// returns a Finding per match. It is extracted from detectLinkname purely to
// flatten the comment-group/comment nesting (the directive prefix check stays at
// depth 2). Matching c.Text against the "//go:linkname" prefix is exact: a prose
// comment whose Text starts with "// " never matches.
func linknameInGroup(group *ast.CommentGroup, fset *token.FileSet, path string) []Finding {
	var findings []Finding

	for _, c := range group.List {
		if strings.HasPrefix(c.Text, "//go:linkname") {
			findings = append(findings, Finding{
				Rule: "TS-REFL-003",
				File: path,
				Line: fset.Position(c.Pos()).Line,
				Note: "//go:linkname directive",
			})
		}
	}

	return findings
}
