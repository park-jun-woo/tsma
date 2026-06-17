//ff:func feature=gate type=test
//ff:what Seed 단위테스트: 명시 root(success)·default root(빈 args)·언어판별 실패·Abs 실패 분기를 임시 fixture로 덮는다. 성공 경로는 Item.Key=QualifiedName과 payload 스냅샷(Status/Attempt 제거)을 검증한다.

package tsmagate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSeed_Success(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module seedtest\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	items, err := New().Seed([]string{root})
	if err != nil {
		t.Fatalf("Seed: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected at least one seeded item")
	}
	var found bool
	for _, it := range items {
		var p funcPayload
		if err := it.DecodePayload(&p); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if p.Lang != "go" {
			t.Errorf("Lang = %q, want go", p.Lang)
		}
		// Root is resolved to an absolute path.
		if !filepath.IsAbs(p.Root) {
			t.Errorf("Root = %q, want absolute", p.Root)
		}
		// Runtime fields are stripped at seed time.
		if p.Fn.Status != "" || p.Fn.Attempt != 0 {
			t.Errorf("payload not stripped: Status=%q Attempt=%d", p.Fn.Status, p.Fn.Attempt)
		}
		if it.Key == "pkg.Foo" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected an item keyed pkg.Foo among %d items", len(items))
	}
}

func TestSeed_DefaultRoot(t *testing.T) {
	// Empty args -> root defaults to "." (the false side of the arg guard).
	root := writeGoPkg(t, map[string]string{
		"go.mod":     "module seedtest\n\ngo 1.22\n",
		"pkg/foo.go": "package pkg\n\nfunc Foo() int { return 1 }\n",
	})
	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	items, err := New().Seed(nil)
	if err != nil {
		t.Fatalf("Seed(nil): %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected items seeded from the default root")
	}
}

func TestSeed_DetectError(t *testing.T) {
	// A directory with no language marker -> detect.Detect errors.
	if _, err := New().Seed([]string{t.TempDir()}); err == nil {
		t.Fatal("expected a language-detection error for an empty dir")
	}
}

func TestSeed_AbsError(t *testing.T) {
	// filepath.Abs fails only when os.Getwd() fails; force it by removing the cwd
	// and passing a relative root so Abs must consult the missing cwd.
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(orig)

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.Remove(dir); err != nil {
		t.Skipf("could not remove cwd: %v", err)
	}
	if _, gErr := os.Getwd(); gErr == nil {
		t.Skip("os.Getwd did not fail after removing cwd on this platform")
	}

	if _, err := New().Seed([]string{"relative-root"}); err == nil {
		t.Fatal("expected an error when filepath.Abs fails")
	}
}
