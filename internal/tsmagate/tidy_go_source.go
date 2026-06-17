//ff:func feature=gate type=helper control=sequence
//ff:what tidyGoSource: sanitize로 언랩한 Go 소스에서 미사용 import를 제거하고 gofmt로 재출력한다. LLM(특히 소형 모델) 테스트 생성 실패의 1위가 "X imported and not used" 컴파일 에러라 디스크에 쓰기 전에 stdlib(go/parser·ast·format)만으로 정리한다. 파싱 실패(잘린 출력 등)면 원본 그대로 반환 — 거부하지 않고 하류 tests-must-pass가 잡게 둔다.

package tsmagate

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
)

// tidyGoSource removes unused imports from the (already unwrapped) Go source and
// re-emits it via gofmt, using only the standard library. The dominant real
// failure of LLM-generated tests is "X imported and not used", so pruning before
// the file lands on disk eliminates that compile error. If the source does not
// parse (e.g. truncated output) it is returned unchanged: tidy never rejects,
// it lets the downstream tests-must-pass gate catch a still-broken file.
func tidyGoSource(src string) string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return src
	}
	used := collectUsedPackages(file)
	pruneUnusedImports(file, used)
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return src
	}
	return buf.String()
}
