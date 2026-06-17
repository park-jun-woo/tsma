//ff:func feature=gate type=test lang=python
//ff:what Python loop e2e (Phase005b §5/§7): drives the real cli.NewQuestCmd tree with a stub backend (llm.CallFunc) through scan→loop against a temp pytest fixture, proving the Python D5 pipeline converges non-invasively. Scenario 1: a full-coverage test PASSes one-shot. Scenario 2: a partial (one-branch) test FAILs on branch coverage, then a full test PASSes — the feedback loop converges. Uses the real ast indexer, content matcher, sys.path-injected .tsma/test backing, and coverage.py branch measurement. Skipped when python3/pytest/coverage are absent.
package tsmagate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/llm"
)

// pyToolsAvailable reports whether python3 with pytest and coverage is usable,
// gating the real-subprocess Python loop e2e (no skip-gate cheating: when the
// tools are present, this test actually runs and measures).
func pyToolsAvailable() bool {
	python, err := exec.LookPath("python3")
	if err != nil {
		return false
	}
	return exec.Command(python, "-c", "import pytest, coverage").Run() == nil
}

// pyFixture writes a minimal pytest project with one branching function and
// returns the project root. pyproject.toml makes detect classify it as Python
// and makes the runner choose pytest (so the backing in .tsma/test is collected
// by path, not by an importable module name).
func pyFixture(t *testing.T) string {
	t.Helper()
	return writeGoPkg(t, map[string]string{
		"pyproject.toml": "[tool.pytest.ini_options]\n",
		"src/calc.py": "def classify(n):\n" +
			"    if n > 0:\n" +
			"        return \"pos\"\n" +
			"    return \"nonpos\"\n",
	})
}

const pyClassifyFullTest = "from calc import classify\n\n" +
	"def test_classify():\n" +
	"    assert classify(1) == \"pos\"\n" +
	"    assert classify(-1) == \"nonpos\"\n"

const pyClassifyPartialTest = "from calc import classify\n\n" +
	"def test_classify():\n" +
	"    assert classify(1) == \"pos\"\n"

// Scenario 1: a complete covering test in one shot → PASS on the first try.
func TestPyLoop_FullCoverageOneShot(t *testing.T) {
	if !pyToolsAvailable() {
		t.Skip("python3 with pytest+coverage not available")
	}
	root := pyFixture(t)
	chdirTo(t, root) // runner/coverage resolve relative test paths against cwd
	session, out := sessionPaths(t)
	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		return pyClassifyFullTest, nil
	})
	opts := loopOpts(backend)

	if _, err := runTsma(t, opts, session, out, "scan", root); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := runTsma(t, opts, session, out, "loop")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if calls != 1 {
		t.Fatalf("backend called %d times, want 1 (one-shot PASS)", calls)
	}
	if !strings.Contains(got, "-> PASS") {
		t.Fatalf("loop output missing PASS:\n%s", got)
	}
	// Non-invasive: the canonical test is promoted beside the source, and the
	// .tsma/test backing is swept (source tree otherwise untouched).
	if _, err := os.Stat(filepath.Join(root, "src", "test_calc.py")); err != nil {
		t.Errorf("expected promoted canonical src/test_calc.py after PASS: %v", err)
	}
}

// Scenario 3: a module that runs code at import time (a side-effect). The
// sys.path-injected backing must import the real module — side-effect and all —
// and still measure to PASS, confirming isolation handles import-time execution
// (plan §5 caveat). _register() runs the moment greet's module is loaded.
func TestPyLoop_ImportSideEffectModule(t *testing.T) {
	if !pyToolsAvailable() {
		t.Skip("python3 with pytest+coverage not available")
	}
	root := writeGoPkg(t, map[string]string{
		"pyproject.toml": "[tool.pytest.ini_options]\n",
		"src/greeter.py": "_LOADED = []\n" +
			"_LOADED.append(\"loaded\")  # import side-effect: runs at module load\n\n" +
			"def greet(name):\n" +
			"    if name:\n" +
			"        return \"hi \" + name\n" +
			"    return \"hi\"\n",
	})
	chdirTo(t, root)
	session, out := sessionPaths(t)
	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		return "from greeter import greet\n\n" +
			"def test_greet():\n" +
			"    assert greet(\"a\") == \"hi a\"\n" +
			"    assert greet(\"\") == \"hi\"\n", nil
	})
	opts := loopOpts(backend)

	if _, err := runTsma(t, opts, session, out, "scan", root); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := runTsma(t, opts, session, out, "loop")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if !strings.Contains(got, "-> PASS") {
		t.Fatalf("side-effect module loop did not PASS:\n%s", got)
	}
}

// Scenario 2: a partial test (one branch) then a full test → FAIL then PASS,
// proving branch coverage drives the retry to convergence.
func TestPyLoop_PartialThenFullConverges(t *testing.T) {
	if !pyToolsAvailable() {
		t.Skip("python3 with pytest+coverage not available")
	}
	root := pyFixture(t)
	chdirTo(t, root)
	session, out := sessionPaths(t)
	calls := 0
	backend := llm.CallFunc(func(system, user string) (string, error) {
		calls++
		if calls == 1 {
			return pyClassifyPartialTest, nil // passes but branch coverage < 100
		}
		return pyClassifyFullTest, nil // retry covers the missing branch
	})
	opts := loopOpts(backend)

	if _, err := runTsma(t, opts, session, out, "scan", root); err != nil {
		t.Fatalf("scan: %v", err)
	}
	got, err := runTsma(t, opts, session, out, "loop")
	if err != nil {
		t.Fatalf("loop: %v", err)
	}
	if calls != 2 {
		t.Fatalf("backend called %d times, want 2 (FAIL then PASS)", calls)
	}
	if !strings.Contains(got, "-> FAIL") || !strings.Contains(got, "-> PASS") {
		t.Fatalf("loop output = %q, want a FAIL then a PASS", got)
	}
}
