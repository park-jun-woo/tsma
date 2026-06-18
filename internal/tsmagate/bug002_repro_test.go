//ff:func feature=gate type=test
//ff:what BUG-002 재현·회귀 테스트: 한 Go 소스 파일(math.go)의 두 함수(Add, Sub)를 차례로 promoteBacking으로 정명(math_test.go)에 확정하면 — 수정 전엔 통째 덮어쓰기라 Add 테스트가 소실되고 Sub만 남았다 — 두 함수의 테스트가 함께 보존됨을 확인한다(누적 promote). 제네릭 경로(prepareLoopNative 없는 언어)는 promoteMerged를 직접 호출하는 동일 경로라 함께 커버. loop의 실제 promote 경로(promoteBacking→promoteMerged)를 LLM 없이 결정적으로 재현한다.

package tsmagate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestBug002_TwoFuncsInOneFileCoexist(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":      "module bug002\n\ngo 1.22\n",
		"pkg/math.go": "package pkg\n\nfunc Add(a, b int) int { return a + b }\nfunc Sub(a, b int) int { return a - b }\n",
	})

	// F1 (Add) passes: its backing is promoted to the canonical math_test.go.
	addBacking := filepath.Join(".tsma", "test", "add.go")
	addSrc := "package pkg\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) { _ = Add(1, 2) }\n"
	if err := writeTestFile(root, addBacking, addSrc); err != nil {
		t.Fatal(err)
	}
	pAdd := funcPayload{Lang: "go", Root: root, Fn: model.Function{QualifiedName: "pkg.Add", File: filepath.Join("pkg", "math.go")}}
	mAdd := &measurement{}
	promoteBacking(pAdd, mAdd, addBacking)
	if mAdd.TestFailed {
		t.Fatalf("Add promote failed: %s", mAdd.FailOutput)
	}

	// F2 (Sub) passes next: SAME canonical path (math_test.go). Pre-fix this
	// overwrote Add's test; post-fix it accumulates as a separate marker block.
	subBacking := filepath.Join(".tsma", "test", "sub.go")
	subSrc := "package pkg\n\nimport \"testing\"\n\nfunc TestSub(t *testing.T) { _ = Sub(3, 1) }\n"
	if err := writeTestFile(root, subBacking, subSrc); err != nil {
		t.Fatal(err)
	}
	pSub := funcPayload{Lang: "go", Root: root, Fn: model.Function{QualifiedName: "pkg.Sub", File: filepath.Join("pkg", "math.go")}}
	mSub := &measurement{}
	promoteBacking(pSub, mSub, subBacking)
	if mSub.TestFailed {
		t.Fatalf("Sub promote failed: %s", mSub.FailOutput)
	}

	got, err := os.ReadFile(filepath.Join(root, "pkg", "math_test.go"))
	if err != nil {
		t.Fatalf("canonical math_test.go must exist: %v", err)
	}
	gs := string(got)
	if !strings.Contains(gs, "func TestAdd") {
		t.Errorf("BUG-002: Add test was lost after Sub promote (overwrite regression):\n%s", gs)
	}
	if !strings.Contains(gs, "func TestSub") {
		t.Errorf("Sub test missing:\n%s", gs)
	}
	if !strings.Contains(gs, "tsma:begin fn=pkg.Add") || !strings.Contains(gs, "tsma:begin fn=pkg.Sub") {
		t.Errorf("both function blocks expected:\n%s", gs)
	}
	if got := strings.Count(gs, "import \"testing\""); got != 1 {
		t.Errorf("shared import should dedup to 1, got %d:\n%s", got, gs)
	}
}
