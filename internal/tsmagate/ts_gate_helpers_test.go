//ff:func feature=gate type=test
//ff:what Phase005a unit tests for the remaining TS gate helpers, each invoked by
//name so it attributes: promoteTS (canonical write / no-path / write-error),
//scanSmells (per-language dispatch), scanTSSmells (findings + error-skip),
//tidyTSSource (best-effort identity), prepareLoopNative (go/typescript/other
//dispatch, driving prepareLoopTS end-to-end), plus extra sanitizeTSSource /
//rewriteTSImports / buildLoopTSTestMatch branch cases.
package tsmagate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

func TestPromoteTS_WritesCanonical(t *testing.T) {
	root := writeTSProj(t, map[string]string{
		"src/math.ts": "export function add(a:number,b:number){return a+b;}\n",
	})
	p := funcPayload{Lang: "typescript", Root: root, Fn: model.Function{File: "src/math.ts", Name: "add"}}
	m := &measurement{}
	src := "import { add } from \"./math\";\ntest('add', () => {});\n"
	promoteTS(p, m, src)
	if m.TestFailed {
		t.Fatalf("promoteTS set TestFailed: %s", m.FailOutput)
	}
	data, err := os.ReadFile(filepath.Join(root, "src", "math.test.ts"))
	if err != nil {
		t.Fatalf("canonical not written: %v", err)
	}
	if !strings.Contains(string(data), `from "./math"`) {
		t.Errorf("canonical lost original import: %s", data)
	}
}

func TestPromoteTS_NoCanonicalPath(t *testing.T) {
	// a non-TS source file yields no canonical TS test path.
	p := funcPayload{Lang: "typescript", Root: t.TempDir(), Fn: model.Function{File: "src/math.go", Name: "add"}}
	m := &measurement{}
	promoteTS(p, m, "x")
	if !m.TestFailed || !strings.Contains(m.FailOutput, "canonical") {
		t.Errorf("expected canonical-path failure, got %+v", m)
	}
}

func TestPromoteTS_WriteError(t *testing.T) {
	// Root is a regular file, so writing under it fails.
	f := filepath.Join(t.TempDir(), "notadir")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "typescript", Root: f, Fn: model.Function{File: "src/math.ts", Name: "add"}}
	m := &measurement{}
	promoteTS(p, m, "x")
	if !m.TestFailed || m.FailOutput == "" {
		t.Errorf("expected write failure, got %+v", m)
	}
}

func TestScanSmells_Dispatch(t *testing.T) {
	// unknown language -> nil.
	if got := scanSmells("rust", "/root", []string{"a_test.rs"}); got != nil {
		t.Errorf("scanSmells(rust) = %v, want nil", got)
	}
	// go dispatch executes the Go scanner (no files -> no findings).
	if got := scanSmells("go", t.TempDir(), nil); len(got) != 0 {
		t.Errorf("scanSmells(go, no files) = %v, want none", got)
	}
	// typescript dispatch executes the TS scanner (missing file -> ignored).
	if got := scanSmells("typescript", t.TempDir(), []string{"missing.test.ts"}); len(got) != 0 {
		t.Errorf("scanSmells(typescript, missing) = %v, want none", got)
	}
	// a language with no detector at all -> nil (the default arm).
	if got := scanSmells("ruby", "/root", []string{"a_test.rb"}); got != nil {
		t.Errorf("scanSmells(ruby) = %v, want nil", got)
	}
}

func TestScanTSSmells_ScanErrorSkipped(t *testing.T) {
	// An unresolvable tree-sitter command makes every ScanTS call error; the
	// scanner must skip the file (continue) and report no findings.
	t.Setenv("TSMA_TREE_SITTER", "tsma-bogus-binary-xyz")
	if got := scanTSSmells(t.TempDir(), []string{"a.test.ts"}); got != nil {
		t.Errorf("scanTSSmells with scan errors = %v, want nil", got)
	}
}

func TestScanTSSmells(t *testing.T) {
	if treesitter.ResolveCommand() == "" {
		t.Skip("tree-sitter unavailable")
	}
	// cheats.test.ts holds the escape hatches; the missing file is ignored.
	got := scanTSSmells("../../testdata/typescript", []string{
		"smell/cheats.test.ts",
		"smell/does-not-exist.test.ts",
	})
	if len(got) == 0 {
		t.Fatal("expected findings from cheats.test.ts")
	}
	for _, f := range got {
		if !strings.HasPrefix(f.Rule, "TS-REFL-TS-") {
			t.Errorf("unexpected rule: %+v", f)
		}
	}
}

func TestTidyTSSource_BestEffort(t *testing.T) {
	// prettier is optional; tidyTSSource must always return non-empty TS, never
	// erroring even when prettier is absent (best-effort identity).
	src := "const x=1;\n"
	got := tidyTSSource(src)
	if strings.TrimSpace(got) == "" {
		t.Errorf("tidyTSSource returned empty for %q", src)
	}
}

func TestSanitizeTSSource_NoCloseAndNoFence(t *testing.T) {
	// opening fence with a language tag but no closing fence: everything after the
	// first newline is kept.
	got := sanitizeTSSource("```ts\nconst a = 1;\n")
	if strings.Contains(got, "```") || !strings.Contains(got, "const a = 1;") {
		t.Errorf("open-fence-only = %q", got)
	}
	// a bare fence token with no trailing newline.
	if got := sanitizeTSSource("```"); strings.Contains(got, "```") {
		t.Errorf("bare fence not stripped: %q", got)
	}
	// no fence at all: the source is returned (trimmed) unchanged.
	got = sanitizeTSSource("  const b = 2;  ")
	if !strings.Contains(got, "const b = 2;") {
		t.Errorf("no-fence path lost body: %q", got)
	}
}

func TestRewriteTSImports_AddsDotPrefixAndBareUntouched(t *testing.T) {
	// the module sits under the backing dir, so the computed relative path has no
	// leading dot and must be prefixed with "./".
	got := rewriteTSImports(`import { f } from "./sub/mod";`, "/proj/.tsma/test/here", "/proj/.tsma/test")
	if !strings.Contains(got, `from "./here/sub/mod"`) {
		t.Errorf("dot-prefix path = %q", got)
	}
	// bare/package specifiers are untouched.
	got = rewriteTSImports(`import x from "lodash";`, "/proj/src", "/proj/.tsma/test")
	if !strings.Contains(got, `"lodash"`) {
		t.Errorf("bare specifier rewritten: %q", got)
	}
}

func TestBuildLoopTSTestMatch_WriteError(t *testing.T) {
	// Root is a regular file so writing the backing under it fails.
	f := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Lang: "typescript", Root: f, Fn: model.Function{File: "src/math.ts", Name: "add"}}
	it := &quest.Item{Key: "src.add"}
	if _, _, err := buildLoopTSTestMatch(p, it, "test('x',()=>{});"); err == nil {
		t.Error("expected write error when Root is a file")
	}
}

func TestPrepareLoopTS_Direct(t *testing.T) {
	// drive prepareLoopTS end-to-end over a real temp project (the measure may
	// pass or fail depending on local tooling, but the pipeline must return a
	// measurement and always sweep its .tsma/test backing).
	root := writeTSProj(t, map[string]string{
		"src/math.ts": "export function add(a, b) { return a + b; }\n",
	})
	p := funcPayload{Lang: "typescript", Root: root, Fn: model.Function{File: "src/math.ts", Name: "add", QualifiedName: "src.add"}}
	raw := "import { add } from \"./math\";\nit('throws to fail fast', () => { throw new Error('x'); });\n"
	m := prepareLoopTS(&quest.Item{Key: "src.add", Tries: 0}, p, []byte(raw))
	if m == nil || m.FuncName != "src.add" {
		t.Fatalf("prepareLoopTS = %+v", m)
	}
	matches, _ := filepath.Glob(filepath.Join(root, ".tsma", "test", "gen-*.test.ts"))
	if len(matches) != 0 {
		t.Errorf("backing not swept: %v", matches)
	}
}

func TestPrepareLoopNative_Dispatch(t *testing.T) {
	// unsupported language -> (nil, false), Prepare then uses the generic path.
	if m, ok := prepareLoopNative(&quest.Item{Key: "x"}, funcPayload{Lang: "ruby"}, nil); ok || m != nil {
		t.Errorf("prepareLoopNative(ruby) = (%v,%v), want (nil,false)", m, ok)
	}

	// Rust (Phase005e): now has a native in-file pipeline. An empty raw declares no
	// #[test] fn, so prepareLoopRs rejects it before any write — exercises the
	// "rust" dispatch arm returning (non-nil, true).
	if m, ok := prepareLoopNative(&quest.Item{Key: "src.calc.add"}, funcPayload{Lang: "rust", Root: t.TempDir(), Fn: model.Function{File: "src/calc.rs", Name: "add"}}, nil); !ok || m == nil || !m.TestFailed {
		t.Errorf("prepareLoopNative(rust) = (%+v,%v), want non-nil TestFailed,true", m, ok)
	}

	// Go: a non-parseable raw short-circuits prepareLoopGo as truncated (fast, no
	// runner) — exercises the "go" dispatch arm.
	m, ok := prepareLoopNative(&quest.Item{Key: "p.add"}, funcPayload{Lang: "go", Root: t.TempDir(), Fn: model.Function{Name: "add"}}, []byte("not go source"))
	if !ok || m == nil || !m.TestFailed || !m.Truncated {
		t.Fatalf("prepareLoopNative(go, garbage) = (%+v,%v)", m, ok)
	}

	// Python: .tsma as a regular file makes the backing write fail inside
	// buildLoopPyTestMatch — exercises the "python" dispatch arm returning a
	// TestFailed measurement without needing a Python toolchain.
	pyRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(pyRoot, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	pyM, pyOK := prepareLoopNative(&quest.Item{Key: "src.calc.classify"}, funcPayload{Lang: "python", Root: pyRoot, Fn: model.Function{File: "src/calc.py", QualifiedName: "src.classify"}}, []byte("def test_x():\n    pass\n"))
	if !pyOK || pyM == nil || !pyM.TestFailed {
		t.Errorf("prepareLoopNative(python) = (%+v,%v), want non-nil TestFailed,true", pyM, pyOK)
	}

	// TypeScript: drives prepareLoopTS end-to-end over a real temp project. The
	// measure may pass or fail depending on local tooling, but the pipeline must
	// always return a measurement and always sweep its .tsma/test backing.
	root := writeTSProj(t, map[string]string{
		"src/math.ts": "export function add(a, b) { return a + b; }\n",
	})
	p := funcPayload{Lang: "typescript", Root: root, Fn: model.Function{File: "src/math.ts", Name: "add", QualifiedName: "src.add"}}
	raw := "import { add } from \"./math\";\nit('throws to fail fast', () => { throw new Error('x'); });\n"
	tm, ok := prepareLoopNative(&quest.Item{Key: "src.add", Tries: 0}, p, []byte(raw))
	if !ok || tm == nil {
		t.Fatalf("prepareLoopNative(typescript) = (%+v,%v), want non-nil,true", tm, ok)
	}
	if tm.FuncName != "src.add" {
		t.Errorf("measurement FuncName = %q, want src.add", tm.FuncName)
	}
	// backing scratch must always be cleaned up.
	matches, _ := filepath.Glob(filepath.Join(root, ".tsma", "test", "gen-*.test.ts"))
	if len(matches) != 0 {
		t.Errorf("backing not swept: %v", matches)
	}
}
