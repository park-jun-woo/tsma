package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestPrintRemainingTodos_listsPartialsAndUntested(t *testing.T) {
	sess := &model.Session{
		Functions: []model.Function{
			{Name: "Pass1", Status: model.StatusPass, CoveragePct: 100},
			{Name: "Partial1", Status: model.StatusTodo, CoveragePct: 92,
				File: "pkg/a.go", StartLine: 1, EndLine: 5},
			{Name: "Untested1", Status: model.StatusTodo,
				File: "pkg/b.go", StartLine: 1, EndLine: 3},
		},
	}
	out := captureStdout(func() { printRemainingTodos(sess) })

	if !strings.Contains(out, "2 TODO") {
		t.Errorf("expected 2 TODO count, got %q", out)
	}
	if !strings.Contains(out, "Partial1") || !strings.Contains(out, "92%") || !strings.Contains(out, "improve coverage") {
		t.Errorf("expected partial line with coverage, got %q", out)
	}
	if !strings.Contains(out, "Untested1") || !strings.Contains(out, "write test") {
		t.Errorf("expected untested line, got %q", out)
	}
	if strings.Contains(out, "Pass1") {
		t.Errorf("PASS functions should not be listed, got %q", out)
	}
}
