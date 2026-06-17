//ff:func feature=smell type=helper control=iteration dimension=1 level=error
//ff:what detectLinkname: TS-REFL-003. 파일의 모든 CommentGroup을 돌며 linknameInGroup으로 `//go:linkname` 컴파일러 지시자를 정확히 탐지한다(c.Text가 "//go:linkname"로 시작). 타 패키지 비공개 심볼을 끌어오는 강한 cheese라 Review로 표면화한다. 지시자가 아닌 일반 주석("// ... //go:linkname ...")이나 문자열 리터럴은 prefix가 다르거나 주석 노드가 아니라 발화하지 않는다(위양성 0).

package smell

import (
	"go/ast"
	"go/token"
)

// detectLinkname finds TS-REFL-003 violations: a //go:linkname compiler
// directive in the test file. It walks each CommentGroup and delegates the
// per-comment scan to linknameInGroup so the comment iteration stays at depth 1.
// A Go directive has no space after // and must start the comment line, so a
// prose comment that merely mentions go:linkname (its Text starts with "// ")
// and a string literal (not a comment node at all) never fire.
func detectLinkname(file *ast.File, fset *token.FileSet, path string) []Finding {
	var findings []Finding

	for _, group := range file.Comments {
		findings = append(findings, linknameInGroup(group, fset, path)...)
	}

	return findings
}
