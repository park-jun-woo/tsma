package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/park-jun-woo/tsma/internal/model"
)

// writeContentAwareModule writes a Go package where scenario_test.go's TestX
// calls Foo by content (not by canonical file name foo_test.go).
func writeContentAwareModule(t *testing.T, dir string) {
	t.Helper()
	srcDir := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo() int { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "scenario_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestScenario(t *testing.T) { _ = Foo() }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDetectTestChange_noTestFile(t *testing.T) {
	dir := t.TempDir()
	// Source file but no test referencing Foo at all.
	srcDir := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo() int { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fn := &model.Function{Name: "Foo", File: "pkg/foo.go"}
	changed, tm := detectTestChange(dir, "go", fn)
	if changed {
		t.Error("expected changed=false when no test references the function")
	}
	if len(tm.Files) != 0 {
		t.Errorf("expected empty match, got %v", tm.Files)
	}
}

func TestDetectTestChange_testFileExists_sameMtime(t *testing.T) {
	dir := t.TempDir()
	writeContentAwareModule(t, dir)

	// Combined mtime equals the stored TestMtime -> not changed.
	testPath := filepath.Join("pkg", "scenario_test.go")
	mtime := getTestMtime(dir, testPath)

	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", TestMtime: mtime}
	changed, tm := detectTestChange(dir, "go", fn)
	if changed {
		t.Error("expected changed=false when mtime matches")
	}
	if len(tm.Files) == 0 {
		t.Error("expected non-empty match")
	}
}

func TestDetectTestChange_testFileExists_differentMtime(t *testing.T) {
	dir := t.TempDir()
	writeContentAwareModule(t, dir)

	fn := &model.Function{Name: "Foo", File: "pkg/foo.go", TestMtime: "2000-01-01T00:00:00Z"}
	changed, tm := detectTestChange(dir, "go", fn)
	if !changed {
		t.Error("expected changed=true when mtime differs")
	}
	if len(tm.Files) == 0 {
		t.Error("expected non-empty match")
	}
	if filepath.Base(tm.Files[0]) != "scenario_test.go" {
		t.Errorf("expected scenario_test.go, got %v", tm.Files[0])
	}
}

// TestDetectTestChange_multiFileMaxMtime verifies the combined mtime is the max
// across multiple matched test files: touching one of them flips changed=true.
func TestDetectTestChange_multiFileMaxMtime(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"),
		[]byte("package pkg\n\nfunc Foo() int { return 1 }\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "a_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestA(t *testing.T) { _ = Foo() }\n"), 0o644)
	os.WriteFile(filepath.Join(srcDir, "b_test.go"),
		[]byte("package pkg\n\nimport \"testing\"\n\nfunc TestB(t *testing.T) { _ = Foo() }\n"), 0o644)

	fn := &model.Function{Name: "Foo", File: "pkg/foo.go"}
	_, tm := detectTestChange(dir, "go", fn)
	if len(tm.Files) != 2 {
		t.Fatalf("expected 2 matched files, got %v", tm.Files)
	}
	mtime := combinedTestMtime(dir, tm.Files)
	fn.TestMtime = mtime
	if changed, _ := detectTestChange(dir, "go", fn); changed {
		t.Error("expected changed=false right after recording combined mtime")
	}

	// Touch one file into the future: combined max advances -> changed.
	future := time.Now().Add(2 * time.Hour)
	os.Chtimes(filepath.Join(srcDir, "b_test.go"), future, future)
	if changed, _ := detectTestChange(dir, "go", fn); !changed {
		t.Error("expected changed=true after touching one of the matched files")
	}
}
