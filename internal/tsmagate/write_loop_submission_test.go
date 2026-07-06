//ff:func feature=gate type=test
//ff:what writeLoopSubmission 직접 호출 테스트: 경로 도출 실패(파생 불가 파일 → error)·기록 실패(정명 경로가 디렉터리 → error)·성공(정명 경로에 생성테스트 기록) 3분기를 임시 fixture로 덮는다.

package tsmagate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestWriteLoopSubmission_TargetPathError(t *testing.T) {
	// A Python payload whose source file is not .py: no test is attributed, no
	// misnamed variant exists, and CanonicalTestPath cannot derive a path ->
	// testTargetPath fails and writeLoopSubmission returns the error (never silent).
	p := funcPayload{
		Lang: "python",
		Root: t.TempDir(),
		Fn:   model.Function{QualifiedName: "app.x", Name: "x", File: filepath.Join("app", "x.txt"), StartLine: 1, EndLine: 1},
	}
	err := writeLoopSubmission(p, []byte("ignored"))
	if err == nil {
		t.Fatal("expected an error when no test path can be derived")
	}
	if !strings.Contains(err.Error(), "cannot derive") {
		t.Fatalf("expected a path-derivation error, got: %v", err)
	}
}

func TestWriteLoopSubmission_WriteError(t *testing.T) {
	// The canonical test path is occupied by a directory, so the merged write
	// fails -> writeLoopSubmission propagates the error for TestFailed surfacing.
	root := t.TempDir()
	appDir := filepath.Join(root, "app")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(appDir, "service.py"), []byte("def do_work():\n    return 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(appDir, "test_service.py"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{
		Lang: "python",
		Root: root,
		Fn:   model.Function{QualifiedName: "app.do_work", Name: "do_work", File: filepath.Join("app", "service.py"), StartLine: 1, EndLine: 2},
	}
	if err := writeLoopSubmission(p, []byte("def test_do_work():\n    assert True\n")); err == nil {
		t.Fatal("expected an error when the canonical test path is a directory")
	}
}

func TestWriteLoopSubmission_Success(t *testing.T) {
	// A brand-new function with no test: the canonical path is derived and the
	// sanitized raw is written there as this function's marker block.
	root := t.TempDir()
	appDir := filepath.Join(root, "app")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(appDir, "service.py"), []byte("def do_work():\n    return 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{
		Lang: "python",
		Root: root,
		Fn:   model.Function{QualifiedName: "app.do_work", Name: "do_work", File: filepath.Join("app", "service.py"), StartLine: 1, EndLine: 2},
	}
	body := "def test_do_work():\n    assert True\n"
	if err := writeLoopSubmission(p, []byte(body)); err != nil {
		t.Fatalf("writeLoopSubmission: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(appDir, "test_service.py"))
	if err != nil {
		t.Fatalf("expected the canonical test file to be written: %v", err)
	}
	if !strings.Contains(string(got), "def test_do_work():") {
		t.Fatalf("written test file is missing the generated body:\n%s", got)
	}
}
