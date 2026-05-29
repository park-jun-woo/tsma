package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestGoFilenameFallback_conventionalFileExists verifies that when the
// conventional <base>_test.go exists on disk the fallback attributes it as a
// single-file TestMatch with TestFuncs left nil.
func TestGoFilenameFallback_conventionalFileExists(t *testing.T) {
	root := t.TempDir()
	pkgDir := filepath.Join("cmd", "app")
	abs := filepath.Join(root, pkgDir)
	if err := os.MkdirAll(abs, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(abs, "agent_cmd.go"), []byte("package app\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(abs, "agent_cmd_test.go"), []byte("package app\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fn := &model.Function{Name: "agentCmd", File: filepath.Join(pkgDir, "agent_cmd.go")}
	tm, ok := goFilenameFallback(root, fn)
	if !ok {
		t.Fatal("expected fallback to match the conventional test file")
	}
	wantFile := filepath.Join(pkgDir, "agent_cmd_test.go")
	if len(tm.Files) != 1 || tm.Files[0] != wantFile {
		t.Fatalf("Files = %v, want [%s]", tm.Files, wantFile)
	}
	if tm.TestFuncs != nil {
		t.Fatalf("TestFuncs = %v, want nil", tm.TestFuncs)
	}
}

// TestGoFilenameFallback_noConventionalFile verifies that when no
// <base>_test.go exists the fallback reports found false and an empty match.
func TestGoFilenameFallback_noConventionalFile(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "lib.go"), []byte("package lib\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fn := &model.Function{Name: "Lib", File: "lib.go"}
	tm, ok := goFilenameFallback(root, fn)
	if ok {
		t.Fatalf("expected no match, got %v", tm)
	}
	if tm.Files != nil || tm.TestFuncs != nil {
		t.Fatalf("expected empty TestMatch, got %v", tm)
	}
}

// TestGoFilenameFallback_nilFunc verifies a nil function reports found false.
func TestGoFilenameFallback_nilFunc(t *testing.T) {
	if tm, ok := goFilenameFallback(t.TempDir(), nil); ok {
		t.Fatalf("expected nil function to be unmatched, got %v", tm)
	}
}
