//ff:func feature=gate type=test
//ff:what Render 단위테스트: payload 디코드 에러·테스트 매칭 성공(test 경로 출력)·미매칭+misnamed 힌트·미매칭+힌트없음 분기를 임시 Go 패키지 fixture로 덮는다.

package tsmagate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/model"
)

// writeGoPkg writes a minimal Go module under root with the given files (keyed by
// path relative to root) and returns root.
func writeGoPkg(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, body := range files {
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// itemWithPayload builds a quest.Item carrying a funcPayload for fn.
func itemWithPayload(t *testing.T, lang, root string, fn model.Function) *quest.Item {
	t.Helper()
	it := &quest.Item{Key: fn.QualifiedName, State: quest.TODO}
	if err := it.SetPayload(funcPayload{Lang: lang, Root: root, Fn: fn}); err != nil {
		t.Fatalf("set payload: %v", err)
	}
	return it
}

func TestRender_DecodeError(t *testing.T) {
	it := &quest.Item{Key: "k", State: quest.TODO, Payload: json.RawMessage("{not json")}
	d := New()
	if _, err := d.Render(nil, it); err == nil {
		t.Fatal("expected a decode error for invalid payload")
	}
}

func TestRender_TestFound(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":          "module rendertest\n\ngo 1.22\n",
		"pkg/foo.go":      "package pkg\n\nfunc Foo() int { return 1 }\n",
		"pkg/foo_test.go": "package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) { _ = Foo() }\n",
	})
	fn := model.Function{QualifiedName: "pkg.Foo", Name: "Foo", File: filepath.Join("pkg", "foo.go"), StartLine: 3, EndLine: 3}
	out, err := New().Render(nil, itemWithPayload(t, "go", root, fn))
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "test: "+filepath.Join("pkg", "foo_test.go")) {
		t.Errorf("output missing matched test path:\n%s", out)
	}
	if strings.Contains(out, "(not found)") {
		t.Errorf("output should not say not found:\n%s", out)
	}
}

func TestRender_NotFoundWithMisnamedHint(t *testing.T) {
	// Canonical pkg/foo_test.go is absent but the misnamed test_foo_test.go
	// exists, so Render surfaces a rename hint.
	root := writeGoPkg(t, map[string]string{
		"go.mod":               "module rendertest\n\ngo 1.22\n",
		"pkg/foo.go":           "package pkg\n\nfunc Foo() int { return 1 }\n",
		"pkg/test_foo_test.go": "package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\n",
	})
	fn := model.Function{QualifiedName: "pkg.Foo", Name: "Foo", File: filepath.Join("pkg", "foo.go"), StartLine: 3, EndLine: 3}
	out, err := New().Render(nil, itemWithPayload(t, "go", root, fn))
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "(not found)") {
		t.Errorf("expected not-found marker:\n%s", out)
	}
	if !strings.Contains(out, "misnamed") {
		t.Errorf("expected a rename hint:\n%s", out)
	}
}

func TestRender_NotFoundNoHint(t *testing.T) {
	// No test of any kind: not found and no rename hint.
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module rendertest\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	fn := model.Function{QualifiedName: "pkg.Foo", Name: "Foo", File: filepath.Join("pkg", "foo.go"), StartLine: 3, EndLine: 3}
	out, err := New().Render(nil, itemWithPayload(t, "go", root, fn))
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "(not found)") {
		t.Errorf("expected not-found marker:\n%s", out)
	}
	if strings.Contains(out, "misnamed") {
		t.Errorf("did not expect a rename hint:\n%s", out)
	}
}
