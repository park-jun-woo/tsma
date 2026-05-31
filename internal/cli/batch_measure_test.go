//ff:test feature=cli
package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestBatchMeasure_DispatchesGo verifies the Go branch of batchMeasure: with
// sess.Lang == "go" it routes through batchMeasureGo, measuring a fully covered
// function to PASS and leaving an untested one TODO, all in one pass.
func TestBatchMeasure_DispatchesGo(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module batchdispatch\n\ngo 1.22\n")
	write("add.go", "package mathx\n\nfunc Add(a, b int) int { return a + b }\n")
	write("mul.go", "package mathx\n\nfunc Mul(a, b int) int { return a * b }\n")
	write("add_test.go", "package mathx\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) { _ = Add(1, 2) }\n")

	sess := &model.Session{
		Lang: "go",
		Functions: []model.Function{
			{Name: "Add", File: "add.go", StartLine: 3, EndLine: 3, Status: model.StatusTodo},
			{Name: "Mul", File: "mul.go", StartLine: 3, EndLine: 3, Status: model.StatusTodo},
		},
	}

	batchMeasure(root, sess)

	add := findFn(sess, "Add")
	mul := findFn(sess, "Mul")
	if add == nil || mul == nil {
		t.Fatalf("missing funcs: add=%v mul=%v", add, mul)
	}
	if add.Status != model.StatusPass {
		t.Errorf("Add: want PASS (Go branch measured it), got %s (%.0f%%)", add.Status, add.CoveragePct)
	}
	if mul.Status != model.StatusTodo {
		t.Errorf("Mul: want TODO (untested), got %s", mul.Status)
	}
}

// TestBatchMeasure_DispatchesOther verifies the non-Go branch of batchMeasure:
// with a non-"go" language it routes through batchMeasureOther, which leaves
// untested / unmeasurable functions TODO without aborting.
func TestBatchMeasure_DispatchesOther(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "lonely.py"), []byte("def lonely():\n    pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	sess := &model.Session{
		Lang: "python",
		Functions: []model.Function{
			{Name: "lonely", File: "lonely.py", StartLine: 1, EndLine: 2, Status: model.StatusTodo},
		},
	}

	batchMeasure(root, sess)

	lonely := findFn(sess, "lonely")
	if lonely == nil {
		t.Fatalf("missing func lonely")
	}
	if lonely.Status != model.StatusTodo {
		t.Errorf("lonely: want TODO (untested, non-Go branch), got %s", lonely.Status)
	}
}
