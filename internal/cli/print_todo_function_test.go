package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestPrintTodoFunction_withTestFile(t *testing.T) {
	fn := &model.Function{
		Name:      "DoSomething",
		File:      "pkg/do.go",
		StartLine: 10,
		EndLine:   20,
	}
	output := captureStdout(func() {
		printTodoFunction(fn, "pkg/do_test.go")
	})

	if !strings.Contains(output, "DoSomething") {
		t.Error("expected function name in output")
	}
	if !strings.Contains(output, "TODO") {
		t.Error("expected TODO in output")
	}
	if !strings.Contains(output, "pkg/do.go:10-20") {
		t.Errorf("expected file:startline-endline in output, got %q", output)
	}
	if !strings.Contains(output, "test: pkg/do_test.go") {
		t.Errorf("expected test file in output, got %q", output)
	}
}

func TestPrintTodoFunction_withoutTestFile(t *testing.T) {
	fn := &model.Function{
		Name:      "NoTest",
		File:      "pkg/no.go",
		StartLine: 1,
		EndLine:   5,
	}
	output := captureStdout(func() {
		printTodoFunction(fn, "")
	})

	if !strings.Contains(output, "(not found)") {
		t.Errorf("expected '(not found)' in output, got %q", output)
	}
}
