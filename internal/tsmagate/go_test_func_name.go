//ff:func feature=gate type=helper control=sequence lang=go
//ff:what goTestFuncName: AST 선언이 패키지 레벨 Go 테스트 함수(func TestXxx(...))면 그 이름과 true를 돌려준다. 리시버가 있거나(메서드) 이름이 Test로 시작하지 않으면 false. parseTestFuncs가 디스크 재매칭 없이 생성 소스에서 -run 대상 테스트명을 직접 뽑을 때 쓴다(중첩 평탄화용 분리).
package tsmagate

import (
	"go/ast"
	"strings"
)

// goTestFuncName reports whether decl is a package-level Go test function
// (func TestXxx(...)) and returns its name. A method (non-nil receiver) or a name
// not prefixed with "Test" is rejected. It exists to keep parseTestFuncs's loop
// at nesting depth ≤2.
func goTestFuncName(decl ast.Decl) (string, bool) {
	fn, ok := decl.(*ast.FuncDecl)
	if !ok || fn.Recv != nil {
		return "", false
	}
	if !strings.HasPrefix(fn.Name.Name, "Test") {
		return "", false
	}
	return fn.Name.Name, true
}
