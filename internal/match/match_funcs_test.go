package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestMatchFuncs_perPackageSeparation verifies that MatchFuncs attributes tests
// using each function's own package index and never bleeds a test from one
// package into a function of another package, even when test func names collide.
func TestMatchFuncs_perPackageSeparation(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/m\n\ngo 1.22\n"), 0o644)

	pkgA := filepath.Join(dir, "a")
	os.MkdirAll(pkgA, 0o755)
	os.WriteFile(filepath.Join(pkgA, "alpha.go"),
		[]byte("package a\n\nfunc Alpha() int { return 1 }\n"), 0o644)
	os.WriteFile(filepath.Join(pkgA, "scenario_test.go"),
		[]byte("package a\n\nimport \"testing\"\n\nfunc TestScenario(t *testing.T) { _ = Alpha() }\n"), 0o644)

	pkgB := filepath.Join(dir, "b")
	os.MkdirAll(pkgB, 0o755)
	os.WriteFile(filepath.Join(pkgB, "beta.go"),
		[]byte("package b\n\nfunc Beta() int { return 2 }\n"), 0o644)
	os.WriteFile(filepath.Join(pkgB, "scenario_test.go"),
		[]byte("package b\n\nimport \"testing\"\n\nfunc TestScenario(t *testing.T) { _ = Beta() }\n"), 0o644)

	fns := []model.Function{
		{Name: "Alpha", File: filepath.Join("a", "alpha.go")},
		{Name: "Beta", File: filepath.Join("b", "beta.go")},
	}

	out := MatchFuncs(dir, fns)

	tmA, ok := out[0]
	if !ok {
		t.Fatal("expected Alpha to be matched")
	}
	if len(tmA.Files) != 1 || filepath.Base(tmA.Files[0]) != "scenario_test.go" {
		t.Errorf("Alpha Files = %v", tmA.Files)
	}
	if filepath.Dir(tmA.Files[0]) != "a" {
		t.Errorf("Alpha matched wrong package: %v", tmA.Files[0])
	}

	tmB, ok := out[1]
	if !ok {
		t.Fatal("expected Beta to be matched")
	}
	if filepath.Dir(tmB.Files[0]) != "b" {
		t.Errorf("Beta matched wrong package: %v", tmB.Files[0])
	}
}

// TestMatchFuncs_unmatchedAbsent verifies functions with no referencing test
// are absent from the result map.
func TestMatchFuncs_unmatchedAbsent(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/m\n\ngo 1.22\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "lib.go"),
		[]byte("package lib\n\nfunc Used() int { return 1 }\n\nfunc Unused() int { return 2 }\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "lib_test.go"),
		[]byte("package lib\n\nimport \"testing\"\n\nfunc TestUsed(t *testing.T) { _ = Used() }\n"), 0o644)

	fns := []model.Function{
		{Name: "Used", File: "lib.go"},
		{Name: "Unused", File: "lib.go"},
	}
	out := MatchFuncs(dir, fns)
	if _, ok := out[0]; !ok {
		t.Error("expected Used to be matched")
	}
	if _, ok := out[1]; ok {
		t.Error("expected Unused to be absent from result")
	}
}

// TestMatchFuncs_emptyInput verifies the empty-slice path returns an empty map.
func TestMatchFuncs_emptyInput(t *testing.T) {
	out := MatchFuncs(t.TempDir(), nil)
	if len(out) != 0 {
		t.Errorf("expected empty result, got %v", out)
	}
}
