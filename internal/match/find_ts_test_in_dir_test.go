package match

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindTSTestInDirTestTS(t *testing.T) {
	dir := t.TempDir()

	srcDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "handler.test.ts"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	testRel, ok := findTSTestInDir(dir, "src", "handler")
	if !ok {
		t.Fatal("expected to find handler.test.ts")
	}
	want := filepath.Join("src", "handler.test.ts")
	if testRel != want {
		t.Errorf("testRel = %q, want %q", testRel, want)
	}
}

func TestFindTSTestInDirSpecTS(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "handler.spec.ts"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	testRel, ok := findTSTestInDir(dir, ".", "handler")
	if !ok {
		t.Fatal("expected to find handler.spec.ts")
	}
	if testRel != "handler.spec.ts" {
		t.Errorf("testRel = %q, want %q", testRel, "handler.spec.ts")
	}
}

func TestFindTSTestInDirTestJS(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "handler.test.js"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	testRel, ok := findTSTestInDir(dir, ".", "handler")
	if !ok {
		t.Fatal("expected to find handler.test.js")
	}
	if testRel != "handler.test.js" {
		t.Errorf("testRel = %q, want %q", testRel, "handler.test.js")
	}
}

func TestFindTSTestInDirNotFound(t *testing.T) {
	dir := t.TempDir()

	_, ok := findTSTestInDir(dir, ".", "handler")
	if ok {
		t.Error("expected not found when no test files exist")
	}
}

func TestFindTSTestInDirPrefersFirstSuffix(t *testing.T) {
	dir := t.TempDir()

	// Create both .test.ts and .spec.ts
	if err := os.WriteFile(filepath.Join(dir, "handler.test.ts"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "handler.spec.ts"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	testRel, ok := findTSTestInDir(dir, ".", "handler")
	if !ok {
		t.Fatal("expected to find test file")
	}
	// .test.ts is first in tsTestSuffixes, so it should be returned
	if testRel != "handler.test.ts" {
		t.Errorf("testRel = %q, want %q (first suffix preferred)", testRel, "handler.test.ts")
	}
}
