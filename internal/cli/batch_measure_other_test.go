//ff:test feature=cli
package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// TestBatchMeasureOther_UntestedAndCheckErrorBothStayTodo exercises both skip
// branches of batchMeasureOther for a non-Go language using Python (a supported
// matcher language):
//   - a function with no attributed test is skipped (stays TODO, not measured);
//   - a function whose test IS attributed (test_greeter.py) but whose
//     checker.Check errors — the test imports a missing module so the coverage
//     run fails — is also skipped (stays TODO), reaching the measure + error-skip
//     path without aborting the scan.
//
// The remaining success line (applyBatchReport on a 100%% report) is delegation
// to applyBatchReport, which is verified directly in apply_batch_report_test.go
// and exercised end-to-end by the Go batch tests; driving it here would require
// the non-Go coverage toolchain to succeed against an absolute project root,
// which the Python checker's fallback path does not support in isolation.
func TestBatchMeasureOther_UntestedAndCheckErrorBothStayTodo(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// greeter.py has a sibling test (matcher attributes it -> measure path); the
	// test imports a missing module so the coverage run errors (the skip-on-error
	// branch). lonely.py has no test (the untested-skip branch).
	write("greeter.py", "def greet():\n    return 'hi'\n")
	write("test_greeter.py", "import does_not_exist\n\ndef test_greet():\n    assert True\n")
	write("lonely.py", "def lonely():\n    pass\n")

	sess := &model.Session{
		Lang: "python",
		Functions: []model.Function{
			{Name: "greet", File: "greeter.py", StartLine: 1, EndLine: 2, Status: model.StatusTodo},
			{Name: "lonely", File: "lonely.py", StartLine: 1, EndLine: 2, Status: model.StatusTodo},
		},
	}

	batchMeasureOther(root, sess)

	for i := range sess.Functions {
		fn := &sess.Functions[i]
		if fn.Status != model.StatusTodo {
			t.Errorf("%s: want TODO (untested or unmeasurable), got %s", fn.Name, fn.Status)
		}
	}
}
