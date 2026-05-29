package match

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRsMatcherInFileTests(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	rel := filepath.Join("src", "lib.rs")
	if err := os.WriteFile(filepath.Join(dir, rel), []byte("pub fn f() {}\n#[cfg(test)]\nmod tests {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &RsMatcher{}
	testFile, found := m.Match(dir, rel)
	if !found {
		t.Fatal("expected match for in-file tests")
	}
	if testFile != rel {
		t.Errorf("testFile = %q, want source file %q", testFile, rel)
	}
}

func TestRsMatcherIntegrationTests(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "src")
	testsDir := filepath.Join(dir, "tests")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(testsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	rel := filepath.Join("src", "lib.rs")
	if err := os.WriteFile(filepath.Join(dir, rel), []byte("pub fn f() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testsDir, "lib.rs"), []byte("// integration\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &RsMatcher{}
	testFile, found := m.Match(dir, rel)
	if !found {
		t.Fatal("expected match for tests/lib.rs")
	}
	want := filepath.Join("tests", "lib.rs")
	if testFile != want {
		t.Errorf("testFile = %q, want %q", testFile, want)
	}
}

func TestRsMatcherNoTest(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	rel := filepath.Join("src", "lib.rs")
	if err := os.WriteFile(filepath.Join(dir, rel), []byte("pub fn f() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &RsMatcher{}
	if _, found := m.Match(dir, rel); found {
		t.Error("expected no match when neither in-file nor tests/ exist")
	}
}
