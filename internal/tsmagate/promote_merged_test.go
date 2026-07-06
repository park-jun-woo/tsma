//ff:func feature=gate type=test
//ff:what promoteMerged 단위테스트: 정명 파일 부재(신규 생성)→존재(누적 머지)의
// read→merge→write 경로와, 부재 이외의 읽기 에러(디렉터리 경로) 전파를 덮는다.
package tsmagate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPromoteMergedAccumulates(t *testing.T) {
	root := t.TempDir()
	rel := filepath.Join("pkg", "x_test.go")
	srcA := "package pkg\n\nimport \"testing\"\n\nfunc TestA(t *testing.T) { _ = 1 }\n"
	if err := promoteMerged(root, rel, srcA, "pkg.A", "go"); err != nil {
		t.Fatalf("first promote (missing canonical): %v", err)
	}
	srcB := "package pkg\n\nimport \"testing\"\n\nfunc TestB(t *testing.T) { _ = 2 }\n"
	if err := promoteMerged(root, rel, srcB, "pkg.B", "go"); err != nil {
		t.Fatalf("second promote (existing canonical): %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatal(err)
	}
	merged := string(data)
	if !strings.Contains(merged, "tsma:begin fn=pkg.A") || !strings.Contains(merged, "tsma:begin fn=pkg.B") {
		t.Errorf("both function blocks must accumulate:\n%s", merged)
	}
}

func TestPromoteMergedReadError(t *testing.T) {
	root := t.TempDir()
	// The canonical path is a directory: ReadFile fails with a non-IsNotExist
	// error, which must propagate (never silently swallowed).
	if err := os.MkdirAll(filepath.Join(root, "adir"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := promoteMerged(root, "adir", "x", "pkg.A", "go"); err == nil {
		t.Error("expected a read error for a directory canonical path")
	}
}
