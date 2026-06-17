//ff:func feature=gate type=test
//ff:what readFile 단위테스트: 파일 내용 반환 성공·읽기 실패("") 분기를 임시 파일로 덮는다.

package tsmagate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile_Success(t *testing.T) {
	dir := t.TempDir()
	abs := filepath.Join(dir, "x.txt")
	const body = "hello\nworld\n"
	if err := os.WriteFile(abs, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := readFile(abs); got != body {
		t.Errorf("readFile = %q, want %q", got, body)
	}
}

func TestReadFile_ReadError(t *testing.T) {
	if got := readFile(filepath.Join(t.TempDir(), "missing.txt")); got != "" {
		t.Errorf("readFile = %q, want empty on read error", got)
	}
}
