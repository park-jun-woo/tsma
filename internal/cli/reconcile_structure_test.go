package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// writeGoSrc writes a Go source file under pkg/ with the given body and returns
// nothing; used to drive the indexer in reconcile tests.
func writeGoSrc(t *testing.T, dir, file, body string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\n\ngo 1.22\n"), 0o644)
	}
	srcDir := filepath.Join(dir, "pkg")
	os.MkdirAll(srcDir, 0o755)
	if err := os.WriteFile(filepath.Join(srcDir, file), []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", file, err)
	}
}

// New source function not in the session is appended as a TODO.
func TestReconcileStructure_addsNewAsTodo(t *testing.T) {
	dir := t.TempDir()
	writeGoSrc(t, dir, "foo.go", "package pkg\n\nfunc Foo() {}\n\nfunc Bar() {}\n")

	sess := &model.Session{
		Lang: "go",
		Functions: []model.Function{
			{QualifiedName: "pkg.Foo", Name: "Foo", File: "pkg/foo.go", Status: model.StatusPass, CoveragePct: 100},
		},
	}

	added, removed, err := reconcileStructure(dir, sess)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if len(added) != 1 || added[0].QualifiedName != "pkg.Bar" {
		t.Fatalf("expected pkg.Bar added, got %+v", added)
	}
	if removed != 0 {
		t.Errorf("expected removed=0, got %d", removed)
	}
	bar := findFn(sess, "Bar")
	if bar == nil || bar.Status != model.StatusTodo {
		t.Errorf("new function must be TODO, got %+v", bar)
	}
	if foo := findFn(sess, "Foo"); foo == nil || foo.Status != model.StatusPass {
		t.Errorf("existing PASS must be preserved, got %+v", foo)
	}
}

// A session function no longer present in source is dropped.
func TestReconcileStructure_removesGone(t *testing.T) {
	dir := t.TempDir()
	writeGoSrc(t, dir, "foo.go", "package pkg\n\nfunc Foo() {}\n")

	sess := &model.Session{
		Lang: "go",
		Functions: []model.Function{
			{QualifiedName: "pkg.Foo", Name: "Foo", File: "pkg/foo.go", Status: model.StatusTodo},
			{QualifiedName: "pkg.Removed", Name: "Removed", File: "pkg/old.go", Status: model.StatusPass},
		},
	}

	added, removed, err := reconcileStructure(dir, sess)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if len(added) != 0 {
		t.Errorf("expected no additions, got %+v", added)
	}
	if removed != 1 {
		t.Errorf("expected removed=1, got %d", removed)
	}
	if findFn(sess, "Removed") != nil {
		t.Error("gone function must be dropped")
	}
	if len(sess.Functions) != 1 {
		t.Errorf("expected 1 function left, got %d", len(sess.Functions))
	}
}

// An existing function that moved/changed keeps its progress while positional
// metadata is refreshed from the new index.
func TestReconcileStructure_preservesProgressRefreshesMetadata(t *testing.T) {
	dir := t.TempDir()
	// Foo now lives lower in the file (line 5) than the stale session says (line 2).
	writeGoSrc(t, dir, "foo.go", "package pkg\n\n// pad\n\nfunc Foo() int {\n\treturn 1\n}\n")

	sess := &model.Session{
		Lang: "go",
		Functions: []model.Function{
			{QualifiedName: "pkg.Foo", Name: "Foo", File: "pkg/old.go", StartLine: 2, EndLine: 2,
				Status: model.StatusDone, CoveragePct: 80, Attempt: 3, TestFiles: []string{"pkg/foo_test.go"}},
		},
	}

	_, _, err := reconcileStructure(dir, sess)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	foo := findFn(sess, "Foo")
	if foo.Status != model.StatusDone || foo.CoveragePct != 80 || foo.Attempt != 3 {
		t.Errorf("progress must be preserved, got %+v", foo)
	}
	if len(foo.TestFiles) != 1 {
		t.Errorf("test files must be preserved, got %+v", foo.TestFiles)
	}
	if foo.File != "pkg/foo.go" || foo.StartLine != 5 {
		t.Errorf("positional metadata must be refreshed, got File=%s StartLine=%d", foo.File, foo.StartLine)
	}
}

// No source change -> no-op (no additions, no removals, status unchanged).
func TestReconcileStructure_noChangeIsNoop(t *testing.T) {
	dir := t.TempDir()
	writeGoSrc(t, dir, "foo.go", "package pkg\n\nfunc Foo() {}\n")

	sess := &model.Session{
		Lang: "go",
		Functions: []model.Function{
			{QualifiedName: "pkg.Foo", Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 3, Status: model.StatusPass},
		},
	}

	added, removed, err := reconcileStructure(dir, sess)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if len(added) != 0 || removed != 0 {
		t.Errorf("expected no-op, got added=%d removed=%d", len(added), removed)
	}
	if sess.Functions[0].Status != model.StatusPass {
		t.Errorf("status must be unchanged, got %s", sess.Functions[0].Status)
	}
}
