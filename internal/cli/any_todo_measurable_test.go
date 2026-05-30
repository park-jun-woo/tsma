package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestAnyTodoMeasurable_trueWhenChanged(t *testing.T) {
	dir := t.TempDir()
	writeGoModule(t, dir) // pkg/foo.go + pkg/foo_test.go covering Foo

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	sess := &model.Session{
		Lang: "go",
		Functions: []model.Function{
			// TestMtime empty -> any real test mtime differs -> changed=true.
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8, Status: model.StatusTodo},
		},
	}
	if !anyTodoMeasurable(dir, "go", sess) {
		t.Error("expected a TODO with a changed test to be measurable")
	}
}

func TestAnyTodoMeasurable_falseWhenUnchangedOrUntested(t *testing.T) {
	dir := t.TempDir()
	testRel := writeGoModule(t, dir)
	mtime := getTestMtime(dir, testRel)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	sess := &model.Session{
		Lang: "go",
		Functions: []model.Function{
			// Stored mtime == current -> unchanged.
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8,
				Status: model.StatusTodo, TestMtime: mtime},
			// PASS functions are ignored entirely.
			{Name: "Bar", File: "pkg/bar.go", Status: model.StatusPass},
		},
	}
	if anyTodoMeasurable(dir, "go", sess) {
		t.Error("expected no measurable TODO (unchanged partial only)")
	}
}

func TestAnyTodoMeasurable_falseWhenNoTest(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	os.WriteFile(filepath.Join(srcDir, "foo.go"), []byte("package pkg\nfunc Foo() {}"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	sess := &model.Session{
		Lang: "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 2, EndLine: 2, Status: model.StatusTodo},
		},
	}
	if anyTodoMeasurable(dir, "go", sess) {
		t.Error("expected untested TODO to be non-measurable")
	}
}
