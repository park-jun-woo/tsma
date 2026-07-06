package tsmagate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/model"
)

func writeRsProj(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestSanitizeRsSource_UnwrapsFence(t *testing.T) {
	raw := "Here is the test:\n```rust\n#[cfg(test)]\nmod tests { #[test] fn t() {} }\n```\nDone."
	got := sanitizeRsSource(raw)
	if strings.Contains(got, "```") || strings.Contains(got, "Here is") {
		t.Errorf("fence/prose not stripped: %q", got)
	}
	if !strings.Contains(got, "#[cfg(test)]") {
		t.Errorf("test body lost: %q", got)
	}
}

func TestSanitizeRsSource_NoFence(t *testing.T) {
	raw := "  #[cfg(test)]\nmod tests { #[test] fn t() {} }  \n"
	got := sanitizeRsSource(raw)
	if !strings.Contains(got, "mod tests") || strings.Contains(got, "```") {
		t.Errorf("no-fence source mangled: %q", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("sanitized source must end with a newline: %q", got)
	}
}

func TestParseRsTestFuncs_BalanceAndNames(t *testing.T) {
	ok := "#[cfg(test)]\nmod tests {\n  #[test]\n  fn a() {}\n  #[tokio::test]\n  async fn b() {}\n}\n"
	names, good := parseRsTestFuncs(ok)
	if !good {
		t.Fatal("balanced source reported truncated")
	}
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Errorf("names = %v, want [a b]", names)
	}
	// Truncated: a missing closing brace.
	if _, good := parseRsTestFuncs("mod tests { #[test] fn a() {"); good {
		t.Error("unbalanced braces should report truncated")
	}
}

func TestInjectRsTestMod_AppendAndReplace(t *testing.T) {
	// Brand-new function: no in-file mod → block is wrapped + appended.
	orig := "pub fn add(a: i32, b: i32) -> i32 { a + b }\n"
	gen := "#[test]\nfn add_works() { assert_eq!(add(1, 2), 3); }"
	out := injectRsTestMod(orig, gen)
	if !strings.Contains(out, "pub fn add") {
		t.Error("production code dropped")
	}
	if !strings.Contains(out, "#[cfg(test)]") || !strings.Contains(out, "add_works") {
		t.Errorf("generated test not wrapped+injected: %s", out)
	}

	// Regenerate: an existing trailing #[cfg(test)] mod is replaced, not duplicated.
	withMod := orig + "\n#[cfg(test)]\nmod tests {\n#[test]\nfn old() {}\n}\n"
	out2 := injectRsTestMod(withMod, "#[cfg(test)]\nmod tests {\n#[test]\nfn fresh() {}\n}")
	if strings.Contains(out2, "fn old()") {
		t.Errorf("stale test mod not replaced: %s", out2)
	}
	if strings.Count(out2, "mod tests") != 1 {
		t.Errorf("duplicate mod tests (compile error): %s", out2)
	}
}

func TestBuildLoopRsTestMatch_InjectsAndBacksUp(t *testing.T) {
	root := writeRsProj(t, map[string]string{
		"src/calc.rs": "pub fn add(a: i32, b: i32) -> i32 { a + b }\n",
	})
	p := funcPayload{Lang: "rust", Root: root, Fn: model.Function{File: "src/calc.rs", Name: "add", QualifiedName: "src.calc.add"}}
	it := &quest.Item{Key: "src.calc.add"}
	gen := "#[cfg(test)]\nmod tests {\n#[test]\nfn add_works() { assert_eq!(add(1,2),3); }\n}\n"

	tm, backingRel, err := buildLoopRsTestMatch(p, it, gen)
	if err != nil {
		t.Fatalf("buildLoopRsTestMatch: %v", err)
	}
	if len(tm.Files) != 1 || tm.Files[0] != "src/calc.rs" {
		t.Fatalf("TestMatch.Files = %v, want [src/calc.rs] (in-file target)", tm.Files)
	}
	// Backing holds the ORIGINAL pre-injection source.
	backing, _ := os.ReadFile(filepath.Join(root, backingRel))
	if strings.Contains(string(backing), "#[cfg(test)]") {
		t.Errorf("backing should be the pristine original, got: %s", backing)
	}
	// The source file itself now carries the injected test (measurement-time mutation).
	src, _ := os.ReadFile(filepath.Join(root, "src/calc.rs"))
	if !strings.Contains(string(src), "add_works") {
		t.Errorf("source not injected: %s", src)
	}
}

func TestFinalizeRsBacking_RollsBackOnNonTerminal(t *testing.T) {
	root := writeRsProj(t, map[string]string{
		"src/calc.rs": "pub fn add(a: i32, b: i32) -> i32 { a + b }\n",
	})
	p := funcPayload{Lang: "rust", Root: root, Fn: model.Function{File: "src/calc.rs", Name: "add"}}
	it := &quest.Item{Key: "src.calc.add", Tries: 0}
	gen := "#[cfg(test)]\nmod tests {\n#[test]\nfn add_works() {}\n}\n"
	_, backingRel, err := buildLoopRsTestMatch(p, it, gen)
	if err != nil {
		t.Fatal(err)
	}

	// A failed measurement → roll back the source, sweep the backing.
	m := &measurement{TestFailed: true}
	finalizeRsBacking(p, it, m, backingRel)

	src, _ := os.ReadFile(filepath.Join(root, "src/calc.rs"))
	if strings.Contains(string(src), "add_works") {
		t.Errorf("source not rolled back: %s", src)
	}
	if _, err := os.Stat(filepath.Join(root, backingRel)); !os.IsNotExist(err) {
		t.Errorf("backing must be swept, stat err = %v", err)
	}
}

func TestFinalizeRsBacking_KeepsInjectionOnTerminalPass(t *testing.T) {
	root := writeRsProj(t, map[string]string{
		"src/calc.rs": "pub fn add(a: i32, b: i32) -> i32 { a + b }\n",
	})
	p := funcPayload{Lang: "rust", Root: root, Fn: model.Function{File: "src/calc.rs", Name: "add"}}
	it := &quest.Item{Key: "src.calc.add", Tries: 0}
	gen := "#[cfg(test)]\nmod tests {\n#[test]\nfn add_works() {}\n}\n"
	_, backingRel, err := buildLoopRsTestMatch(p, it, gen)
	if err != nil {
		t.Fatal(err)
	}

	// A passing, fully-covered measurement → keep the injected source.
	m := &measurement{Report: &coverage.Report{AllCovered: true}}
	finalizeRsBacking(p, it, m, backingRel)

	src, _ := os.ReadFile(filepath.Join(root, "src/calc.rs"))
	if !strings.Contains(string(src), "add_works") {
		t.Errorf("terminal pass must keep the in-file test: %s", src)
	}
	if _, err := os.Stat(filepath.Join(root, backingRel)); !os.IsNotExist(err) {
		t.Errorf("backing must be swept even on materialize, stat err = %v", err)
	}
}

func TestPrepareLoopRs_TruncatedAndTestless(t *testing.T) {
	root := writeRsProj(t, map[string]string{"src/calc.rs": "pub fn add() {}\n"})
	p := funcPayload{Lang: "rust", Root: root, Fn: model.Function{File: "src/calc.rs", Name: "add", QualifiedName: "src.calc.add"}}
	it := &quest.Item{Key: "src.calc.add"}

	// Truncated (unbalanced) → TestFailed+Truncated, source untouched.
	m := prepareLoopRs(it, p, []byte("mod tests { #[test] fn t() {"))
	if !m.TestFailed || !m.Truncated {
		t.Errorf("truncated block: TestFailed=%v Truncated=%v, want both true", m.TestFailed, m.Truncated)
	}
	// No #[test] fn → rejected before any write.
	m2 := prepareLoopRs(it, p, []byte("#[cfg(test)]\nmod tests { fn helper() {} }\n"))
	if !m2.TestFailed || m2.Truncated {
		t.Errorf("testless block: TestFailed=%v Truncated=%v, want TestFailed only", m2.TestFailed, m2.Truncated)
	}
	src, _ := os.ReadFile(filepath.Join(root, "src/calc.rs"))
	if strings.Contains(string(src), "mod tests") {
		t.Errorf("rejected generation must not mutate source: %s", src)
	}
}

func TestPrepareLoopRs_CargoAbsentRestoresSource(t *testing.T) {
	// cargo is absent in this environment (the documented e2e skip-gate): the
	// runner errors, measureLoop sets TestFailed, and the deferred finalize rolls
	// the source back — the injection/rollback logic is exercised end-to-end
	// without cargo.
	if _, err := exec.LookPath("cargo"); err == nil {
		t.Skip("cargo present: this test asserts the cargo-absent rollback path")
	}
	root := writeRsProj(t, map[string]string{
		"src/calc.rs": "pub fn add(a: i32, b: i32) -> i32 { a + b }\n",
		"Cargo.toml":  "[package]\nname=\"x\"\nversion=\"0.1.0\"\n",
	})
	p := funcPayload{Lang: "rust", Root: root, Fn: model.Function{File: "src/calc.rs", Name: "add", QualifiedName: "src.calc.add"}}
	it := &quest.Item{Key: "src.calc.add"}
	gen := "#[cfg(test)]\nmod tests {\n#[test]\nfn add_works() { assert_eq!(add(1,2),3); }\n}\n"

	m := prepareLoopRs(it, p, []byte(gen))
	if !m.TestFailed {
		t.Errorf("cargo absent must surface TestFailed, got %+v", m)
	}
	src, _ := os.ReadFile(filepath.Join(root, "src/calc.rs"))
	if strings.Contains(string(src), "add_works") {
		t.Errorf("source must be rolled back after a non-terminal measurement: %s", src)
	}
}
