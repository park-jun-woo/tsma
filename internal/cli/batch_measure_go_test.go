//ff:test feature=cli
package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// setupBatchRepo writes a tiny Go module in one package: Add (fully covered ->
// PASS), Sub (partial -> measured TODO), Mul (untested -> TODO).
func setupBatchRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	write := func(name, body string) {
		if werr := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); werr != nil {
			t.Fatal(werr)
		}
	}
	write("go.mod", "module batchtest\n\ngo 1.22\n")
	write("add.go", "package mathx\n\nfunc Add(a, b int) int { return a + b }\n")
	write("sub.go", "package mathx\n\nfunc Sub(a, b int) int {\n\tif a > b {\n\t\treturn a - b\n\t}\n\treturn b - a\n}\n")
	write("mul.go", "package mathx\n\nfunc Mul(a, b int) int { return a * b }\n")
	write("add_test.go", "package mathx\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) { _ = Add(1, 2) }\n")
	write("sub_test.go", "package mathx\n\nimport \"testing\"\n\nfunc TestSub(t *testing.T) { _ = Sub(3, 1) }\n")
	return root
}

func findFn(sess *model.Session, name string) *model.Function {
	for i := range sess.Functions {
		if sess.Functions[i].Name == name {
			return &sess.Functions[i]
		}
	}
	return nil
}

// TestBatchMeasureGo_MeasuresAllInOnePass: a single analyze (one go test per
// package) yields PASS for the 100% function, a measured TODO for the partial,
// and a plain TODO for the untested function — never an auto-DONE.
func TestBatchMeasureGo_MeasuresAllInOnePass(t *testing.T) {
	root := setupBatchRepo(t)
	sess, err := analyzeProject(root)
	if err != nil {
		t.Fatalf("analyzeProject: %v", err)
	}

	if !sess.FirstPassDone {
		t.Fatalf("expected FirstPassDone after analysis (batch should finish measuring)")
	}

	add := findFn(sess, "Add")
	sub := findFn(sess, "Sub")
	mul := findFn(sess, "Mul")
	if add == nil || sub == nil || mul == nil {
		t.Fatalf("missing funcs: add=%v sub=%v mul=%v", add, sub, mul)
	}

	if add.Status != model.StatusPass {
		t.Errorf("Add: want PASS, got %s (%.0f%%)", add.Status, add.CoveragePct)
	}
	if sub.Status != model.StatusTodo {
		t.Errorf("Sub: want TODO (partial), got %s", sub.Status)
	}
	if sub.CoveragePct <= 0 || sub.CoveragePct >= 100 {
		t.Errorf("Sub: want measured partial (0<pct<100), got %.1f", sub.CoveragePct)
	}
	if sub.Attempt != 1 {
		t.Errorf("Sub: want Attempt=1 (batch measured once, no auto-DONE), got %d", sub.Attempt)
	}
	if sub.Status == model.StatusDone {
		t.Errorf("Sub: batch must never auto-promote a partial to DONE")
	}
	if mul.Status != model.StatusTodo {
		t.Errorf("Mul: want TODO (untested), got %s", mul.Status)
	}
}

// TestBatchMeasureGo_CompileFailureSkipsPackageNotScan: a compile-failing
// package is skipped (its functions stay TODO) while a healthy package is still
// measured to PASS — the scan is never aborted.
func TestBatchMeasureGo_CompileFailureSkipsPackageNotScan(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		if werr := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); werr != nil {
			t.Fatal(werr)
		}
	}
	mkdir := func(name string) {
		if derr := os.MkdirAll(filepath.Join(root, name), 0o755); derr != nil {
			t.Fatal(derr)
		}
	}
	write("go.mod", "module batchfail\n\ngo 1.22\n")
	mkdir("ok")
	write("ok/add.go", "package okp\n\nfunc Add(a, b int) int { return a + b }\n")
	write("ok/add_test.go", "package okp\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) { _ = Add(1, 2) }\n")
	mkdir("bad")
	write("bad/x.go", "package badp\n\nfunc Broken() int { return 1 }\n")
	write("bad/x_test.go", "package badp\n\nimport \"testing\"\n\nfunc TestBroken(t *testing.T) { _ = Broken(); undefinedSymbol() }\n")

	sess, err := analyzeProject(root)
	if err != nil {
		t.Fatalf("analyzeProject must not abort on a compile-failing package: %v", err)
	}

	add := findFn(sess, "Add")
	broken := findFn(sess, "Broken")
	if add == nil || broken == nil {
		t.Fatalf("missing funcs: add=%v broken=%v", add, broken)
	}
	if add.Status != model.StatusPass {
		t.Errorf("Add (good pkg) must still be measured PASS, got %s", add.Status)
	}
	if broken.Status != model.StatusTodo {
		t.Errorf("Broken (failing pkg) must stay TODO, got %s", broken.Status)
	}
}
