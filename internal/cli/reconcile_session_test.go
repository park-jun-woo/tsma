package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// measure=false: a new source function is added as a TODO but never measured.
func TestReconcileSession_structureOnlyLeavesNewUnmeasured(t *testing.T) {
	dir := t.TempDir()
	writeGoSrc(t, dir, "foo.go", "package pkg\n\nfunc Foo() {}\n\nfunc New() {}\n")

	sess := &model.Session{
		Lang: "go",
		Functions: []model.Function{
			{QualifiedName: "pkg.Foo", Name: "Foo", File: "pkg/foo.go", Status: model.StatusPass, CoveragePct: 100},
		},
	}

	added, removed, err := reconcileSession(dir, sess, false)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if added != 1 || removed != 0 {
		t.Fatalf("expected +1/-0, got +%d/-%d", added, removed)
	}
	n := findFn(sess, "New")
	if n == nil || n.Status != model.StatusTodo {
		t.Fatalf("new func must be TODO, got %+v", n)
	}
	if n.Attempt != 0 || n.CoveragePct != 0 {
		t.Errorf("measure=false must not measure: got Attempt=%d Pct=%v", n.Attempt, n.CoveragePct)
	}
	if sess.Summary.Total != 2 || sess.Summary.Todo != 1 || sess.Summary.Pass != 1 {
		t.Errorf("summary recalculated wrong: %+v", sess.Summary)
	}
}

// BUG-004 regression: a fully-PASS stale session that has gained a new source
// function must NOT report "All functions complete!" on `tsma next`; the new
// function must surface as a TODO.
func TestRunNext_staleSessionSurfacesNewFunction(t *testing.T) {
	dir := t.TempDir()
	writeGoSrc(t, dir, "foo.go", "package pkg\n\nfunc Foo() {}\n")
	// A newly-extracted helper that the stale session never indexed (no test).
	writeGoSrc(t, dir, "helper.go", "package pkg\n\nfunc Helper() {}\n")

	// Stale session: only Foo, already PASS -> would falsely be "complete".
	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{QualifiedName: "pkg.Foo", Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 3,
				Status: model.StatusPass, CoveragePct: 100},
		},
		CurrentIndex: 0,
		Summary:      model.Summary{Total: 1, Pass: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdoutAndStderr(func() {
		if err := runNext(nil, nil); err != nil {
			t.Fatalf("runNext: %v", err)
		}
	})

	if strings.Contains(out, "All functions complete!") {
		t.Errorf("must not report completion while an untracked function exists, got %q", out)
	}
	if !strings.Contains(out, "Helper") {
		t.Errorf("new function Helper must be surfaced, got %q", out)
	}

	loaded := loadSessionForTest(t, dir)
	if loaded.Summary.Total != 2 {
		t.Errorf("expected 2 functions after reconcile, got %d", loaded.Summary.Total)
	}
	if h := findFn(loaded, "Helper"); h == nil || h.Status != model.StatusTodo {
		t.Errorf("Helper must be tracked as TODO, got %+v", h)
	}
	if f := findFn(loaded, "Foo"); f == nil || f.Status != model.StatusPass {
		t.Errorf("existing Foo PASS must be preserved, got %+v", f)
	}
}
