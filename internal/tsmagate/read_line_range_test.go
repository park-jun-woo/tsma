//ff:func feature=gate type=test
//ff:what readLineRange 단위테스트: 정상 범위·읽기 실패·start<1·end<start·start>len·end>len(끝 클램프) 분기를 임시 파일로 덮는다.

package tsmagate

import (
	"os"
	"path/filepath"
	"testing"
)

func writeLines(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	abs := filepath.Join(dir, "lines.txt")
	if err := os.WriteFile(abs, []byte("l1\nl2\nl3\nl4\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return abs
}

func TestReadLineRange_Normal(t *testing.T) {
	abs := writeLines(t)
	if got := readLineRange(abs, 2, 3); got != "l2\nl3" {
		t.Errorf("readLineRange = %q, want %q", got, "l2\nl3")
	}
}

func TestReadLineRange_ReadError(t *testing.T) {
	if got := readLineRange(filepath.Join(t.TempDir(), "missing.txt"), 1, 1); got != "" {
		t.Errorf("readLineRange = %q, want empty on read error", got)
	}
}

func TestReadLineRange_StartBelowOne(t *testing.T) {
	abs := writeLines(t)
	if got := readLineRange(abs, 0, 2); got != "" {
		t.Errorf("readLineRange = %q, want empty when start < 1", got)
	}
}

func TestReadLineRange_EndBeforeStart(t *testing.T) {
	abs := writeLines(t)
	if got := readLineRange(abs, 3, 2); got != "" {
		t.Errorf("readLineRange = %q, want empty when end < start", got)
	}
}

func TestReadLineRange_StartBeyondEOF(t *testing.T) {
	abs := writeLines(t)
	if got := readLineRange(abs, 99, 100); got != "" {
		t.Errorf("readLineRange = %q, want empty when start past EOF", got)
	}
}

func TestReadLineRange_EndClampedToEOF(t *testing.T) {
	abs := writeLines(t)
	// File has 4 lines + a trailing empty element after the final newline.
	// Asking past the end clamps to the last line rather than panicking.
	got := readLineRange(abs, 3, 99)
	if got == "" {
		t.Fatalf("readLineRange returned empty, want clamped tail")
	}
	if got[:2] != "l3" {
		t.Errorf("readLineRange = %q, want it to start at l3", got)
	}
}
