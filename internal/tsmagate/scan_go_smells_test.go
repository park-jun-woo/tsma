//ff:func feature=gate type=test
//ff:what scanGoSmells 단위테스트: 깨끗한 _test.go(unsafe finding → append 분기)와 파싱 불가 파일(ScanGo 에러 → continue 분기)을 함께 넘겨 두 루프 분기를 모두 덮는다.

package tsmagate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanGoSmells_AppendAndContinue(t *testing.T) {
	dir := t.TempDir()

	good := "good_test.go"
	goodSrc := "package p\n\nimport \"unsafe\"\n\nfunc TestX() {\n\tvar x int\n\t_ = unsafe.Pointer(&x)\n}\n"
	if err := os.WriteFile(filepath.Join(dir, good), []byte(goodSrc), 0o644); err != nil {
		t.Fatalf("write good: %v", err)
	}

	bad := "bad_test.go"
	if err := os.WriteFile(filepath.Join(dir, bad), []byte("package p\nfunc ("), 0o644); err != nil {
		t.Fatalf("write bad: %v", err)
	}

	// good yields findings (append branch); bad fails to parse (continue branch).
	got := scanGoSmells(dir, []string{good, bad})
	if len(got) == 0 {
		t.Fatal("expected at least one finding from the unsafe import in good_test.go")
	}
	for _, f := range got {
		if f.Rule != "TS-REFL-001" {
			t.Errorf("Rule = %q, want TS-REFL-001", f.Rule)
		}
	}
}
