//ff:func feature=gate type=test
//ff:what goPackageName 단위테스트: package 라인 추출 성공·package 라인 없음(루프 소진 후 "")·읽기 실패("") 분기를 임시 파일로 덮는다.

package tsmagate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoPackageName_Found(t *testing.T) {
	dir := t.TempDir()
	abs := filepath.Join(dir, "foo.go")
	if err := os.WriteFile(abs, []byte("// a comment\n\npackage pkg\n\nfunc Foo() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := goPackageName(abs); got != "pkg" {
		t.Errorf("goPackageName = %q, want pkg", got)
	}
}

func TestGoPackageName_NoPackageLine(t *testing.T) {
	dir := t.TempDir()
	abs := filepath.Join(dir, "nopkg.go")
	if err := os.WriteFile(abs, []byte("// just a comment\nvar x = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := goPackageName(abs); got != "" {
		t.Errorf("goPackageName = %q, want empty when no package line", got)
	}
}

func TestGoPackageName_ReadError(t *testing.T) {
	if got := goPackageName(filepath.Join(t.TempDir(), "does-not-exist.go")); got != "" {
		t.Errorf("goPackageName = %q, want empty on read error", got)
	}
}
