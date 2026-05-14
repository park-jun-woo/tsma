package match

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPyMatcherMatchSameDir(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "handler.py"), []byte("def handler(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "test_handler.py"), []byte("def test_handler(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &PyMatcher{}
	testFile, found := m.Match(dir, "handler.py")
	if !found {
		t.Fatal("expected to find test file in same directory")
	}
	if testFile != "test_handler.py" {
		t.Errorf("testFile = %q, want %q", testFile, "test_handler.py")
	}
}

func TestPyMatcherMatchTestsSubdir(t *testing.T) {
	dir := t.TempDir()

	srcDir := filepath.Join(dir, "app")
	testsDir := filepath.Join(dir, "app", "tests")
	if err := os.MkdirAll(testsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "auth.py"), []byte("def login(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testsDir, "test_auth.py"), []byte("def test_login(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &PyMatcher{}
	testFile, found := m.Match(dir, "app/auth.py")
	if !found {
		t.Fatal("expected to find test file in tests/ subdir")
	}
	want := filepath.Join("app", "tests", "test_auth.py")
	if testFile != want {
		t.Errorf("testFile = %q, want %q", testFile, want)
	}
}

func TestPyMatcherMatchNotFound(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "handler.py"), []byte("def handler(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &PyMatcher{}
	_, found := m.Match(dir, "handler.py")
	if found {
		t.Error("expected no match when no test file exists")
	}
}

func TestPyMatcherMatchPrefersSameDir(t *testing.T) {
	dir := t.TempDir()

	testsDir := filepath.Join(dir, "tests")
	if err := os.MkdirAll(testsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "handler.py"), []byte("def handler(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "test_handler.py"), []byte("def test_handler(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testsDir, "test_handler.py"), []byte("def test_handler(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &PyMatcher{}
	testFile, found := m.Match(dir, "handler.py")
	if !found {
		t.Fatal("expected to find test file")
	}
	// Same dir should be preferred (searched first)
	if testFile != "test_handler.py" {
		t.Errorf("testFile = %q, want %q (same dir preferred)", testFile, "test_handler.py")
	}
}
