//ff:func feature=smell type=helper control=sequence level=error
//ff:what ScanGo: 한 _test.go 파일을 go/ast로 1회 파싱(주석 포함)해 세 detector(go_unsafe/go_reflect_dynamic/go_linkname)를 돌려 []Finding을 합쳐 돌려준다. 파싱 에러는 호출자(prepare)가 무시하도록 그대로 반환(빌드 깨진 테스트는 tests-must-pass가 잡음). substring이 아니라 AST 노드 기준이라 주석·문자열 리터럴 속 키워드는 발화하지 않는다(위양성 0).

package smell

import (
	"go/parser"
	"go/token"
)

// ScanGo parses a single _test.go file with go/ast (comments included so the
// //go:linkname directive is visible) and collects escape-hatch findings from
// the three Go detectors. It returns (nil, err) on a parse error so the caller
// can ignore a broken test file (tests-must-pass owns build failures). It never
// reports the legitimate reflect.DeepEqual/TypeOf/ValueOf idioms (false-positive
// zero) because every detector matches AST nodes, not source substrings.
func ScanGo(path string) ([]Finding, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var findings []Finding
	findings = append(findings, detectUnsafe(file, fset, path)...)
	findings = append(findings, detectReflectDynamic(file, fset, path)...)
	findings = append(findings, detectLinkname(file, fset, path)...)
	return findings, nil
}
