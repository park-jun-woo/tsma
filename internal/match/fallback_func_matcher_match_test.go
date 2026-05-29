package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestFallbackFuncMatcherWrapsMatch verifies the fallback adapter returns the
// exact same test file as the wrapped legacy Matcher.Match, with nil TestFuncs.
func TestFallbackFuncMatcherWrapsMatch(t *testing.T) {
	root := t.TempDir()
	srcRel := filepath.Join("app", "service.py")
	testRel := filepath.Join("app", "test_service.py")
	if err := os.MkdirAll(filepath.Join(root, "app"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, srcRel), []byte("# src\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, testRel), []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	inner := &PyMatcher{}
	wantFile, wantOk := inner.Match(root, srcRel)
	if !wantOk {
		t.Fatal("setup: expected legacy PyMatcher to find the test")
	}

	fm := NewFuncMatcher("python")
	fn := &model.Function{Name: "do_work", File: srcRel}
	tm, ok := fm.MatchFunc(root, fn)
	if !ok {
		t.Fatal("expected fallback to find the test")
	}
	if len(tm.Files) != 1 || tm.Files[0] != wantFile {
		t.Fatalf("Files = %v, want [%s] (same as legacy Match)", tm.Files, wantFile)
	}
	if tm.TestFuncs != nil {
		t.Fatalf("TestFuncs = %v, want nil (runner resolves from file)", tm.TestFuncs)
	}
}

// TestFallbackFuncMatcherNotFound verifies the adapter reports found false when
// the legacy Matcher finds nothing, returning an empty match.
func TestFallbackFuncMatcherNotFound(t *testing.T) {
	root := t.TempDir()
	srcRel := filepath.Join("app", "service.py")
	if err := os.MkdirAll(filepath.Join(root, "app"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, srcRel), []byte("# src\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fm := NewFuncMatcher("python")
	fn := &model.Function{Name: "do_work", File: srcRel}
	if tm, ok := fm.MatchFunc(root, fn); ok {
		t.Fatalf("expected no match, got %v", tm)
	}
}

// TestFallbackFuncMatcherNilFunc verifies a nil function yields found false.
func TestFallbackFuncMatcherNilFunc(t *testing.T) {
	fm := &fallbackFuncMatcher{inner: &PyMatcher{}}
	if _, ok := fm.MatchFunc(t.TempDir(), nil); ok {
		t.Fatal("expected nil function to be unmatched")
	}
}
