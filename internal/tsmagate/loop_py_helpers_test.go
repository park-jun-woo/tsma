//ff:func feature=gate type=test lang=python
//ff:what Phase005b D5 unit tests for the Python loop helpers, each called by name
// so the content matcher attributes coverage to it: pyBackingSlug (identifier
// sanitization), injectPySyspath (sys.path header), runPyFormatter/tidyPySource/
// sanitizePySource (best-effort fence-unwrap + format), buildLoopPyTestMatch
// (isolated backing under .tsma/test + write-error), promotePy/finalizePyBacking
// (canonical written from ORIGINAL src only on a terminal pass; backing always
// swept), and prepareLoopPy (build-failure surfaced as TestFailed; full pipeline
// gated on python tools). The pytest-backed measureLoop runs only when the tools
// are present.
package tsmagate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/model"
)

func TestPyBackingSlug_SanitizesToIdentifier(t *testing.T) {
	tests := []struct{ in, want string }{
		{"src.calc.classify", "src_calc_classify"},
		{"Foo-Bar 1", "Foo_Bar_1"},
		{"already_ok_9", "already_ok_9"},
		{"", ""},
		{"a/b(c)[d]", "a_b_c__d_"},
	}
	for _, tc := range tests {
		if got := pyBackingSlug(tc.in); got != tc.want {
			t.Errorf("pyBackingSlug(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestInjectPySyspath_PrependsHeader(t *testing.T) {
	got := injectPySyspath("def test_x():\n    pass\n", "/proj/src")
	if !strings.HasPrefix(got, "import sys as _tsma_sys\n") {
		t.Errorf("missing sys import header: %q", got)
	}
	if !strings.Contains(got, `_tsma_sys.path.insert(0, "/proj/src")`) {
		t.Errorf("missing sys.path.insert with quoted dir: %q", got)
	}
	if !strings.HasSuffix(got, "def test_x():\n    pass\n") {
		t.Errorf("original body not preserved at the end: %q", got)
	}
}

func TestRunPyFormatter_Branches(t *testing.T) {
	// cat echoes stdin verbatim (success, non-empty output).
	if got := runPyFormatter("x = 1\n", "cat"); got != "x = 1\n" {
		t.Errorf("cat passthrough = %q, want %q", got, "x = 1\n")
	}
	// `true` exits 0 with no output → empty result → src unchanged.
	if got := runPyFormatter("orig = 2\n", "true"); got != "orig = 2\n" {
		t.Errorf("empty output should fall back to src, got %q", got)
	}
	// A missing tool fails to run → src unchanged.
	if got := runPyFormatter("kept = 3\n", "definitely-not-a-formatter-xyz"); got != "kept = 3\n" {
		t.Errorf("missing tool should fall back to src, got %q", got)
	}
}

func TestTidyPySource_BestEffort(t *testing.T) {
	// isort/black are optional; tidy must never lose the code body.
	got := tidyPySource("def foo():\n    return 1\n")
	if !strings.Contains(got, "def foo") {
		t.Errorf("tidyPySource dropped the body: %q", got)
	}
}

func TestSanitizePySource_UnwrapsFences(t *testing.T) {
	t.Run("fence with language tag and closing fence", func(t *testing.T) {
		raw := "Here is the test:\n```python\ndef foo():\n    return 1\n```\nDone."
		got := sanitizePySource(raw)
		if strings.Contains(got, "```") || strings.Contains(got, "Here is") || strings.Contains(got, "Done.") {
			t.Errorf("fence/prose not stripped: %q", got)
		}
		if !strings.Contains(got, "def foo") {
			t.Errorf("body lost: %q", got)
		}
	})

	t.Run("fence without closing fence (j<0)", func(t *testing.T) {
		raw := "```python\ndef bar():\n    return 2"
		got := sanitizePySource(raw)
		if strings.Contains(got, "```") || strings.Contains(got, "python\n") {
			t.Errorf("opening fence/lang tag not stripped: %q", got)
		}
		if !strings.Contains(got, "def bar") {
			t.Errorf("body lost: %q", got)
		}
	})

	t.Run("no fence trims whitespace", func(t *testing.T) {
		got := sanitizePySource("\n\ndef baz():\n    return 3\n\n")
		if strings.Contains(got, "```") {
			t.Errorf("unexpected fence: %q", got)
		}
		if !strings.Contains(got, "def baz") {
			t.Errorf("body lost: %q", got)
		}
	})

	t.Run("bare fence with no newline (nl<0)", func(t *testing.T) {
		// Exercises the nl<0 path: nothing follows the opening ``` on its line.
		_ = sanitizePySource("```")
	})
}

func TestBuildLoopPyTestMatch_IsolatesAndInjects(t *testing.T) {
	root := writeGoPkg(t, map[string]string{
		"src/calc.py": "def classify(n):\n    return n\n",
	})
	p := funcPayload{Lang: "python", Root: root, Fn: model.Function{File: "src/calc.py", Name: "classify"}}
	it := &quest.Item{Key: "src.calc.classify"}
	src := "from calc import classify\n\ndef test_classify():\n    assert classify(1) == 1\n"

	tm, backingRel, err := buildLoopPyTestMatch(p, it, src)
	if err != nil {
		t.Fatalf("buildLoopPyTestMatch: %v", err)
	}
	if len(tm.Files) != 1 || tm.Files[0] != backingRel {
		t.Fatalf("TestMatch.Files = %v, want [%s]", tm.Files, backingRel)
	}
	wantPrefix := filepath.Join(".tsma", "test", "gen_src_calc_classify.py")
	if backingRel != wantPrefix {
		t.Errorf("backing path = %s, want %s", backingRel, wantPrefix)
	}
	data, err := os.ReadFile(filepath.Join(root, backingRel))
	if err != nil {
		t.Fatalf("backing not written: %v", err)
	}
	if !strings.Contains(string(data), "_tsma_sys.path.insert(0,") {
		t.Errorf("backing missing the sys.path injection header: %s", data)
	}
	if !strings.Contains(string(data), "from calc import classify") {
		t.Errorf("backing missing the original test body: %s", data)
	}
	// Source tree must stay clean (no test_calc.py beside the source).
	if _, err := os.Stat(filepath.Join(root, "src", "test_calc.py")); !os.IsNotExist(err) {
		t.Errorf("source tree must stay clean during measurement, stat err = %v", err)
	}
}

func TestBuildLoopPyTestMatch_WriteError(t *testing.T) {
	root := t.TempDir()
	// Make .tsma a regular file so MkdirAll(.tsma/test) fails.
	if err := os.WriteFile(filepath.Join(root, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "python", Root: root, Fn: model.Function{File: "src/calc.py", Name: "classify"}}
	it := &quest.Item{Key: "src.calc.classify"}
	if _, _, err := buildLoopPyTestMatch(p, it, "def test_x():\n    pass\n"); err == nil {
		t.Fatal("expected a write error when .tsma is a file")
	}
}

func TestPromotePy_Branches(t *testing.T) {
	t.Run("non-py source yields no canonical path", func(t *testing.T) {
		p := funcPayload{Lang: "python", Root: t.TempDir(), Fn: model.Function{File: "src/notes.txt"}}
		m := &measurement{}
		promotePy(p, m, "whatever")
		if !m.TestFailed || !strings.Contains(m.FailOutput, "cannot derive canonical") {
			t.Fatalf("expected TestFailed with derive error, got %+v", m)
		}
	})

	t.Run("success writes the canonical test beside the source", func(t *testing.T) {
		root := writeGoPkg(t, map[string]string{"src/calc.py": "def classify(n):\n    return n\n"})
		p := funcPayload{Lang: "python", Root: root, Fn: model.Function{File: "src/calc.py"}}
		m := &measurement{}
		promotePy(p, m, "from calc import classify\n")
		if m.TestFailed {
			t.Fatalf("unexpected failure: %s", m.FailOutput)
		}
		data, err := os.ReadFile(filepath.Join(root, "src", "test_calc.py"))
		if err != nil {
			t.Fatalf("canonical not written: %v", err)
		}
		if !strings.Contains(string(data), "from calc import classify") {
			t.Errorf("canonical must carry the ORIGINAL src: %s", data)
		}
	})

	t.Run("write error surfaces as TestFailed", func(t *testing.T) {
		root := t.TempDir()
		// Make src a regular file so MkdirAll(src) fails for src/test_calc.py.
		if err := os.WriteFile(filepath.Join(root, "src"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		p := funcPayload{Lang: "python", Root: root, Fn: model.Function{File: "src/calc.py"}}
		m := &measurement{}
		promotePy(p, m, "content")
		if !m.TestFailed || m.FailOutput == "" {
			t.Fatalf("expected TestFailed with output, got %+v", m)
		}
	})
}

func TestFinalizePyBacking_PromotesOnPassSweepsAlways(t *testing.T) {
	root := writeGoPkg(t, map[string]string{"src/calc.py": "def classify(n):\n    return n\n"})
	p := funcPayload{Lang: "python", Root: root, Fn: model.Function{File: "src/calc.py"}}
	it := &quest.Item{Key: "src.calc.classify", Tries: 0}
	backingRel := filepath.Join(".tsma", "test", "gen_src_calc_classify.py")
	if err := writeTestFile(root, backingRel, "injected backing"); err != nil {
		t.Fatal(err)
	}
	originalSrc := "from calc import classify\n\ndef test_classify():\n    assert classify(1) == 1\n"

	m := &measurement{Report: &coverage.Report{AllCovered: true}}
	finalizePyBacking(p, it, m, originalSrc, backingRel)

	data, err := os.ReadFile(filepath.Join(root, "src", "test_calc.py"))
	if err != nil {
		t.Fatalf("canonical not promoted on pass: %v", err)
	}
	if !strings.Contains(string(data), "from calc import classify") {
		t.Errorf("canonical must carry the ORIGINAL src (no sys.path header): %s", data)
	}
	if _, err := os.Stat(filepath.Join(root, backingRel)); !os.IsNotExist(err) {
		t.Errorf("backing must be swept, stat err = %v", err)
	}
}

func TestFinalizePyBacking_FailDoesNotTouchSource(t *testing.T) {
	root := writeGoPkg(t, map[string]string{"src/calc.py": "def classify(n):\n    return n\n"})
	p := funcPayload{Lang: "python", Root: root, Fn: model.Function{File: "src/calc.py"}}
	it := &quest.Item{Key: "src.calc.classify", Tries: 0}
	backingRel := filepath.Join(".tsma", "test", "gen_src_calc_classify.py")
	if err := writeTestFile(root, backingRel, "injected backing"); err != nil {
		t.Fatal(err)
	}

	m := &measurement{TestFailed: true}
	finalizePyBacking(p, it, m, "x", backingRel)

	if _, err := os.Stat(filepath.Join(root, "src", "test_calc.py")); !os.IsNotExist(err) {
		t.Errorf("a failed loop must not write into the source tree, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, backingRel)); !os.IsNotExist(err) {
		t.Errorf("backing must be swept even on failure, stat err = %v", err)
	}
}

func TestPrepareLoopPy_BuildFailureIsTestFailed(t *testing.T) {
	root := t.TempDir()
	// .tsma as a file → backing write fails inside buildLoopPyTestMatch.
	if err := os.WriteFile(filepath.Join(root, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "python", Root: root, Fn: model.Function{File: "src/calc.py", QualifiedName: "src.classify"}}
	it := &quest.Item{Key: "src.calc.classify"}
	m := prepareLoopPy(it, p, []byte("def test_x():\n    pass\n"))
	if !m.TestFailed || m.FailOutput == "" {
		t.Fatalf("expected TestFailed with output on a write failure, got %+v", m)
	}
}

func TestPrepareLoopPy_FullPipelinePasses(t *testing.T) {
	if !pyToolsAvailable() {
		t.Skip("python3 with pytest+coverage not available")
	}
	root := pyFixture(t)
	chdirTo(t, root)
	p := funcPayload{
		Lang: "python",
		Root: root,
		Fn:   model.Function{File: "src/calc.py", Name: "classify", QualifiedName: "src.classify", StartLine: 1, EndLine: 4},
	}
	it := &quest.Item{Key: "src.calc.classify", Tries: 0}
	m := prepareLoopPy(it, p, []byte(pyClassifyFullTest))
	if m.TestFailed {
		t.Fatalf("full-coverage test should pass, got TestFailed: %s", m.FailOutput)
	}
	if m.Report == nil {
		t.Fatalf("expected a coverage report on pass, got %+v", m)
	}
	if !m.Report.AllCovered {
		t.Errorf("expected full branch coverage, report = %+v", m.Report)
	}
	// Terminal pass (100%) promotes the canonical test beside the source.
	if _, err := os.Stat(filepath.Join(root, "src", "test_calc.py")); err != nil {
		t.Errorf("expected promoted canonical src/test_calc.py after PASS: %v", err)
	}
	// The .tsma/test backing must be swept.
	if _, err := os.Stat(filepath.Join(root, ".tsma", "test", "gen_src_calc_classify.py")); !os.IsNotExist(err) {
		t.Errorf("backing must be swept after measurement, stat err = %v", err)
	}
}
