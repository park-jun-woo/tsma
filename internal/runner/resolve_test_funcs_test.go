package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/match"
)

// TestResolveTestFuncs_explicit returns content-aware TestFuncs verbatim.
func TestResolveTestFuncs_explicit(t *testing.T) {
	m := match.TestMatch{Files: []string{"a_test.go"}, TestFuncs: []string{"TestX", "TestY"}}
	got, err := ResolveTestFuncs(m)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "TestX" || got[1] != "TestY" {
		t.Errorf("got %v, want [TestX TestY]", got)
	}
}

// TestResolveTestFuncs_extractUnion parses files when TestFuncs is nil and
// returns the deduplicated union across files.
func TestResolveTestFuncs_extractUnion(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a_test.go")
	f2 := filepath.Join(dir, "b_test.go")
	os.WriteFile(f1, []byte("package p\nimport \"testing\"\nfunc TestA(t *testing.T){}\nfunc TestShared(t *testing.T){}\n"), 0o644)
	os.WriteFile(f2, []byte("package p\nimport \"testing\"\nfunc TestB(t *testing.T){}\nfunc TestShared(t *testing.T){}\n"), 0o644)

	m := match.TestMatch{Files: []string{f1, f2}, TestFuncs: nil}
	got, err := ResolveTestFuncs(m)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"TestA": true, "TestB": true, "TestShared": true}
	if len(got) != 3 {
		t.Fatalf("got %v, want 3 unique funcs", got)
	}
	for _, name := range got {
		if !want[name] {
			t.Errorf("unexpected func %q in %v", name, got)
		}
	}
}

func TestResolveTestFuncs_missingFileError(t *testing.T) {
	m := match.TestMatch{Files: []string{"/nonexistent/x_test.go"}, TestFuncs: nil}
	if _, err := ResolveTestFuncs(m); err == nil {
		t.Error("expected error for missing file")
	}
}

// TestResolveTestFuncs_absError covers the filepath.Abs error branch (line 24):
// Abs fails only when os.Getwd() fails, forced by removing the cwd, with a
// relative match file so Abs must consult the missing cwd.
func TestResolveTestFuncs_absError(t *testing.T) {
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

	m := match.TestMatch{Files: []string{"rel_test.go"}, TestFuncs: nil}
	if _, err := ResolveTestFuncs(m); err == nil {
		t.Fatal("expected error when filepath.Abs fails")
	}
}
