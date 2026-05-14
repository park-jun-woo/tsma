package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintFuncList_mixed(t *testing.T) {
	functions := []model.Function{
		{Name: "FuncA", Status: model.StatusPass, CoveragePct: 100},
		{Name: "FuncBB", Status: model.StatusDone, CoveragePct: 75},
		{Name: "FuncCCC", Status: model.StatusTodo},
	}
	maxName := maxFuncNameLen(functions)

	output := captureStdout(func() {
		printFuncList(functions, maxName)
	})

	if !strings.Contains(output, "FuncA") {
		t.Error("expected FuncA in output")
	}
	if !strings.Contains(output, "PASS") {
		t.Error("expected PASS in output")
	}
	if !strings.Contains(output, "FuncBB") {
		t.Error("expected FuncBB in output")
	}
	if !strings.Contains(output, "DONE") {
		t.Error("expected DONE in output")
	}
	if !strings.Contains(output, "FuncCCC") {
		t.Error("expected FuncCCC in output")
	}
	if !strings.Contains(output, "TODO") {
		t.Error("expected TODO in output")
	}
}

func TestPrintFuncList_empty(t *testing.T) {
	output := captureStdout(func() {
		printFuncList([]model.Function{}, 0)
	})
	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

func TestPrintFuncList_passCoverage(t *testing.T) {
	functions := []model.Function{
		{Name: "Fn", Status: model.StatusPass, CoveragePct: 100},
	}
	output := captureStdout(func() {
		printFuncList(functions, 2)
	})
	if !strings.Contains(output, "100%") {
		t.Errorf("expected 100%% in output, got %q", output)
	}
}
