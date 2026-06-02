package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// rescan with no session errors out.
func TestRunRescan_noSession(t *testing.T) {
	dir := t.TempDir()

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	if err := runRescan(nil, nil); err == nil {
		t.Fatal("expected error when no session exists")
	}
}

// rescan syncs a new function in while preserving existing progress, without a
// reset.
func TestRunRescan_addsNewPreservesProgress(t *testing.T) {
	dir := t.TempDir()
	writeGoSrc(t, dir, "foo.go", "package pkg\n\nfunc Foo() {}\n\nfunc Added() {}\n")

	sess := &model.Session{
		Project: dir,
		Lang:    "go",
		Functions: []model.Function{
			{QualifiedName: "pkg.Foo", Name: "Foo", File: "pkg/foo.go", Status: model.StatusDone, CoveragePct: 75, Attempt: 2},
		},
		Summary: model.Summary{Total: 1, Done: 1},
	}
	writeSession(t, dir, sess)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	out := captureStdout(func() {
		if err := runRescan(nil, nil); err != nil {
			t.Fatalf("rescan: %v", err)
		}
	})
	if !strings.Contains(out, "+1 new") {
		t.Errorf("expected '+1 new' in output, got %q", out)
	}

	loaded := loadSessionForTest(t, dir)
	if loaded.Summary.Total != 2 {
		t.Errorf("expected 2 functions, got %d", loaded.Summary.Total)
	}
	if foo := findFn(loaded, "Foo"); foo == nil || foo.Status != model.StatusDone || foo.Attempt != 2 {
		t.Errorf("existing DONE progress must be preserved, got %+v", foo)
	}
	if added := findFn(loaded, "Added"); added == nil || added.Status != model.StatusTodo {
		t.Errorf("new function must be tracked as TODO, got %+v", added)
	}
}
