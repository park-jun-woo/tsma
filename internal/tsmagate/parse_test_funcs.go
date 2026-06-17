//ff:func feature=gate type=helper control=iteration dimension=1 lang=go
//ff:what parseTestFuncs: sanitize+tidy를 거친 생성 소스를 go/parser로 파싱해 테스트 함수명을 직접 추출한다(디스크 재매칭 우회 — C2). 파싱 성공이면 (이름들, true), 실패면 (nil, false)을 돌려준다. false는 곧 잘림/미완성(C3)이라 Prepare가 측정을 생략하고 truncated 피드백을 낸다.
package tsmagate

import (
	"go/parser"
	"go/token"
)

// parseTestFuncs parses the (already sanitized+tidied) generated source and
// returns the names of its package-level test functions. A second return of
// false means the source did not parse — i.e. the model's output was truncated —
// which the loop turns into the C3 truncated feedback instead of measuring. It
// deliberately does not touch disk: the test funcs come straight from the AST,
// bypassing the disk re-match that would miss a brand-new test (C2).
func parseTestFuncs(src string) ([]string, bool) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, false
	}
	var funcs []string
	for _, decl := range file.Decls {
		if name, ok := goTestFuncName(decl); ok {
			funcs = append(funcs, name)
		}
	}
	return funcs, true
}
