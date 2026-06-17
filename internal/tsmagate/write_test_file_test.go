//ff:func feature=gate type=test
//ff:what writeTestFile 단위테스트: 디렉터리 생성 후 정상 기록·MkdirAll 실패(부모 경로가 파일)·WriteFile 실패(대상이 디렉터리) 분기를 임시 트리로 덮는다.

package tsmagate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteTestFile_Success(t *testing.T) {
	root := t.TempDir()
	const src = "package pkg\n\nfunc TestX(t *testing.T) {}\n"
	if err := writeTestFile(root, filepath.Join("pkg", "x_test.go"), src); err != nil {
		t.Fatalf("writeTestFile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(root, "pkg", "x_test.go"))
	if err != nil {
		t.Fatalf("read written file: %v", err)
	}
	if string(got) != src {
		t.Errorf("written content = %q, want %q", got, src)
	}
}

func TestWriteTestFile_MkdirError(t *testing.T) {
	root := t.TempDir()
	// "pkg" is a regular file, so creating "pkg/" as a directory must fail.
	if err := os.WriteFile(filepath.Join(root, "pkg"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeTestFile(root, filepath.Join("pkg", "x_test.go"), "src"); err == nil {
		t.Fatal("expected an error when the parent path is a file")
	}
}

func TestWriteTestFile_WriteError(t *testing.T) {
	root := t.TempDir()
	// The target itself is a directory, so WriteFile to it must fail.
	if err := os.MkdirAll(filepath.Join(root, "adir"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := writeTestFile(root, "adir", "src"); err == nil {
		t.Fatal("expected an error when the target is a directory")
	}
}
