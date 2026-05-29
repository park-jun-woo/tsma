package match

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// writeFuncMatcherPkg writes files into a fresh temp package dir under
// internal/gen and returns the project root and the package-relative dir.
func writeFuncMatcherPkg(t *testing.T, files map[string]string) (root, pkgDir string) {
	t.Helper()
	root = t.TempDir()
	pkgDir = filepath.Join("internal", "gen")
	abs := filepath.Join(root, pkgDir)
	if err := os.MkdirAll(abs, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(abs, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root, pkgDir
}

func sortedCopy(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}

// TestGoFuncMatcherExactAttribution verifies GenerateBytes-style exact
// attribution: GenerateBytes resolves only to its own test, never bleeding into
// Generate (the decisive content-aware case).
func TestGoFuncMatcherExactAttribution(t *testing.T) {
	root, pkgDir := writeFuncMatcherPkg(t, map[string]string{
		"generate_write_file_test.go": `package gen
import "testing"
func TestWriteFile(t *testing.T) { Generate("a", "b") }
`,
		"generate_escrow_test.go": `package gen
import "testing"
func TestEscrow(t *testing.T) { GenerateBytes("a", nil) }
`,
	})

	m := &GoFuncMatcher{}
	fn := &model.Function{Name: "GenerateBytes", File: filepath.Join(pkgDir, "generate.go")}
	tm, ok := m.MatchFunc(root, fn)
	if !ok {
		t.Fatal("expected GenerateBytes to be attributed")
	}
	wantFile := filepath.Join(pkgDir, "generate_escrow_test.go")
	if len(tm.Files) != 1 || tm.Files[0] != wantFile {
		t.Fatalf("Files = %v, want [%s]", tm.Files, wantFile)
	}
	if len(tm.TestFuncs) != 1 || tm.TestFuncs[0] != "TestEscrow" {
		t.Fatalf("TestFuncs = %v, want [TestEscrow]", tm.TestFuncs)
	}
}

// TestGoFuncMatcherOneToMany verifies a single source function attributed to
// multiple test files and test functions, deduplicated.
func TestGoFuncMatcherOneToMany(t *testing.T) {
	root, pkgDir := writeFuncMatcherPkg(t, map[string]string{
		"parse_happy_test.go": `package gen
import "testing"
func TestParseHappy(t *testing.T) { Parse("x") }
`,
		"parse_error_test.go": `package gen
import "testing"
func TestParseError(t *testing.T) { Parse("") }
func TestParseEmpty(t *testing.T) { Parse("  ") }
`,
	})

	m := &GoFuncMatcher{}
	fn := &model.Function{Name: "Parse", File: filepath.Join(pkgDir, "parse.go")}
	tm, ok := m.MatchFunc(root, fn)
	if !ok {
		t.Fatal("expected Parse to be attributed")
	}
	gotFiles := sortedCopy(tm.Files)
	wantFiles := []string{
		filepath.Join(pkgDir, "parse_error_test.go"),
		filepath.Join(pkgDir, "parse_happy_test.go"),
	}
	if len(gotFiles) != 2 || gotFiles[0] != wantFiles[0] || gotFiles[1] != wantFiles[1] {
		t.Fatalf("Files = %v, want %v", gotFiles, wantFiles)
	}
	gotFuncs := sortedCopy(tm.TestFuncs)
	wantFuncs := []string{"TestParseEmpty", "TestParseError", "TestParseHappy"}
	if len(gotFuncs) != 3 || gotFuncs[0] != wantFuncs[0] || gotFuncs[1] != wantFuncs[1] || gotFuncs[2] != wantFuncs[2] {
		t.Fatalf("TestFuncs = %v, want %v", gotFuncs, wantFuncs)
	}
}

// TestGoFuncMatcherMethod verifies that a method call x.Fetch() is attributed
// to a function whose bare Name is the method name.
func TestGoFuncMatcherMethod(t *testing.T) {
	root, pkgDir := writeFuncMatcherPkg(t, map[string]string{
		"client_test.go": `package gen
import "testing"
func TestClient(t *testing.T) {
	c := NewClient()
	c.Fetch("url")
}
`,
	})

	m := &GoFuncMatcher{}
	fn := &model.Function{Name: "Fetch", File: filepath.Join(pkgDir, "client.go")}
	tm, ok := m.MatchFunc(root, fn)
	if !ok {
		t.Fatal("expected Fetch to be attributed")
	}
	if len(tm.TestFuncs) != 1 || tm.TestFuncs[0] != "TestClient" {
		t.Fatalf("TestFuncs = %v, want [TestClient]", tm.TestFuncs)
	}
}

// TestGoFuncMatcherUntestedNotFound verifies an untested function yields found
// false and an empty match (no file-name fallback).
func TestGoFuncMatcherUntestedNotFound(t *testing.T) {
	root, pkgDir := writeFuncMatcherPkg(t, map[string]string{
		"only_test.go": `package gen
import "testing"
func TestOnly(t *testing.T) { Used() }
`,
	})

	m := &GoFuncMatcher{}
	fn := &model.Function{Name: "Unused", File: filepath.Join(pkgDir, "unused.go")}
	tm, ok := m.MatchFunc(root, fn)
	if ok {
		t.Fatalf("expected Unused to be unmatched, got %v", tm)
	}
	if tm.Files != nil || tm.TestFuncs != nil {
		t.Fatalf("expected empty TestMatch, got %v", tm)
	}
}

// TestGoFuncMatcherNoFileNameFallback verifies that even when a canonically
// named test file exists, attribution is purely content-aware: a function whose
// name no test references is not matched despite the conventional file name.
func TestGoFuncMatcherNoFileNameFallback(t *testing.T) {
	root, pkgDir := writeFuncMatcherPkg(t, map[string]string{
		"widget_test.go": `package gen
import "testing"
func TestWidget(t *testing.T) { Other() }
`,
	})

	m := &GoFuncMatcher{}
	// Source file widget.go has a conventional widget_test.go, but no test
	// references the function Widget, so content-aware must not match it.
	fn := &model.Function{Name: "Widget", File: filepath.Join(pkgDir, "widget.go")}
	if tm, ok := m.MatchFunc(root, fn); ok {
		t.Fatalf("content-aware must not fall back to file name, got %v", tm)
	}
}

// TestGoFuncMatcherNilFunc verifies a nil function yields found false.
func TestGoFuncMatcherNilFunc(t *testing.T) {
	m := &GoFuncMatcher{}
	if _, ok := m.MatchFunc(t.TempDir(), nil); ok {
		t.Fatal("expected nil function to be unmatched")
	}
}

// TestGoFuncMatcherMissingDir verifies a non-existent package dir yields found
// false rather than panicking.
func TestGoFuncMatcherMissingDir(t *testing.T) {
	m := &GoFuncMatcher{}
	fn := &model.Function{Name: "X", File: filepath.Join("nope", "x.go")}
	if _, ok := m.MatchFunc(t.TempDir(), fn); ok {
		t.Fatal("expected missing dir to be unmatched")
	}
}
