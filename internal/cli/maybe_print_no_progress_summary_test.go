package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestMaybePrintNoProgressSummary_printsWhenStuck(t *testing.T) {
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
	out := captureStdout(func() { maybePrintNoProgressSummary(dir, sess) })
	if !strings.Contains(out, "TODO function(s) remaining") {
		t.Errorf("expected summary when no TODO is measurable, got %q", out)
	}
}

func TestMaybePrintNoProgressSummary_silentWhenMeasurable(t *testing.T) {
	dir := t.TempDir()
	writeGoModule(t, dir)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	sess := &model.Session{
		Lang: "go",
		Functions: []model.Function{
			{Name: "Foo", File: "pkg/foo.go", StartLine: 3, EndLine: 8, Status: model.StatusTodo},
		},
	}
	out := captureStdout(func() { maybePrintNoProgressSummary(dir, sess) })
	if out != "" {
		t.Errorf("expected no output while progress is possible, got %q", out)
	}
}
