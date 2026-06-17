package tsmagate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// TestStripRsCfgTestMod covers the no-marker (unchanged) and trailing-mod
// (stripped back to the marker) branches.
func TestStripRsCfgTestMod(t *testing.T) {
	noMod := "pub fn f() -> i32 { 1 }\n"
	if got := stripRsCfgTestMod(noMod); got != noMod {
		t.Errorf("no marker: stripRsCfgTestMod = %q, want unchanged", got)
	}

	withMod := "pub fn f() -> i32 { 1 }\n\n#[cfg(test)]\nmod tests {\n#[test]\nfn t() {}\n}\n"
	got := stripRsCfgTestMod(withMod)
	if strings.Contains(got, "#[cfg(test)]") || strings.Contains(got, "mod tests") {
		t.Errorf("trailing mod not stripped: %q", got)
	}
	if got != "pub fn f() -> i32 { 1 }\n" {
		t.Errorf("stripRsCfgTestMod = %q, want production code + newline", got)
	}
}

// TestTidyRsSource covers the rustfmt-absent passthrough (the only deterministic
// branch in this environment); when rustfmt is present the output is still valid
// Rust, so we only assert it preserves the test marker.
func TestTidyRsSource(t *testing.T) {
	src := "#[cfg(test)]\nmod tests {\n#[test]\nfn t() { assert_eq!(1, 1); }\n}\n"
	got := tidyRsSource(src)
	if _, err := exec.LookPath("rustfmt"); err != nil {
		if got != src {
			t.Errorf("rustfmt absent: tidyRsSource changed src: %q", got)
		}
		return
	}
	if !strings.Contains(got, "#[cfg(test)]") {
		t.Errorf("rustfmt present: test marker lost: %q", got)
	}
}

// TestRestoreRsSource covers the successful copy-back and the read-error
// (missing backing) branch.
func TestRestoreRsSource(t *testing.T) {
	root := writeRsProj(t, map[string]string{
		"src/calc.rs":             "INJECTED CONTENT\n",
		".tsma/test/orig-calc.rs": "pub fn add() {}\n",
	})
	p := funcPayload{Root: root, Fn: model.Function{File: "src/calc.rs"}}

	if err := restoreRsSource(p, filepath.Join(".tsma", "test", "orig-calc.rs")); err != nil {
		t.Fatalf("restoreRsSource: %v", err)
	}
	src, _ := os.ReadFile(filepath.Join(root, "src/calc.rs"))
	if string(src) != "pub fn add() {}\n" {
		t.Errorf("source not restored from backing: %q", src)
	}

	// missing backing -> error, source left untouched.
	if err := restoreRsSource(p, filepath.Join(".tsma", "test", "orig-missing.rs")); err == nil {
		t.Error("missing backing: restoreRsSource err = nil, want error")
	}
	src2, _ := os.ReadFile(filepath.Join(root, "src/calc.rs"))
	if string(src2) != "pub fn add() {}\n" {
		t.Errorf("source mutated despite restore error: %q", src2)
	}
}

// TestScanRsSmells covers the success-append branch (tree-sitter over cheese.rs
// yields escape-hatch findings) and the error-continue branch (a missing file),
// plus the no-tree-sitter empty result.
func TestScanRsSmells(t *testing.T) {
	// error path: a missing file is skipped (continue), no panic, no findings.
	if got := scanRsSmells("../..", []string{"testdata/rust/src/does-not-exist.rs"}); got != nil {
		t.Errorf("missing file: scanRsSmells = %+v, want nil", got)
	}

	if treesitter.ResolveCommand() == "" || treesitter.ResolveGrammar("rust") == "" {
		t.Skip("tree-sitter CLI + rust grammar not available")
	}
	// success path: cheese.rs has unsafe/transmute/ptr in its #[cfg(test)] module.
	got := scanRsSmells("../..", []string{"testdata/rust/src/cheese.rs", "testdata/rust/src/does-not-exist.rs"})
	if len(got) == 0 {
		t.Error("expected findings from cheese.rs test scope")
	}
}

// TestBuildLoopRsTestMatch_ReadError covers the original-read failure branch
// (the source file does not exist).
func TestBuildLoopRsTestMatch_ReadError(t *testing.T) {
	root := t.TempDir()
	p := funcPayload{Root: root, Fn: model.Function{File: "src/missing.rs", Name: "add"}}
	it := &quest.Item{Key: "src.missing.add"}
	if _, _, err := buildLoopRsTestMatch(p, it, "gen"); err == nil {
		t.Error("missing source: buildLoopRsTestMatch err = nil, want error")
	}
}

// TestBuildLoopRsTestMatch_BackingWriteError covers the backing-write failure
// branch: a regular file at .tsma makes MkdirAll(.tsma/test) fail.
func TestBuildLoopRsTestMatch_BackingWriteError(t *testing.T) {
	root := writeRsProj(t, map[string]string{
		"src/calc.rs": "pub fn add() {}\n",
	})
	// place a regular file where the .tsma directory must be created.
	if err := os.WriteFile(filepath.Join(root, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	p := funcPayload{Root: root, Fn: model.Function{File: "src/calc.rs", Name: "add"}}
	it := &quest.Item{Key: "src.calc.add"}
	if _, _, err := buildLoopRsTestMatch(p, it, "gen"); err == nil {
		t.Error("blocked backing dir: buildLoopRsTestMatch err = nil, want error")
	}
}

// TestBuildLoopRsTestMatch_SourceWriteError covers the injected-source write
// failure branch: the source is read-only, so the backup read succeeds but
// re-writing the source with the injected mod fails. Skipped under root (perms
// are ignored there).
func TestBuildLoopRsTestMatch_SourceWriteError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root: file permissions are ignored")
	}
	root := writeRsProj(t, map[string]string{
		"src/calc.rs": "pub fn add() {}\n",
	})
	absSrc := filepath.Join(root, "src/calc.rs")
	if err := os.Chmod(absSrc, 0o444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(absSrc, 0o644) })

	p := funcPayload{Root: root, Fn: model.Function{File: "src/calc.rs", Name: "add"}}
	it := &quest.Item{Key: "src.calc.add"}
	if _, _, err := buildLoopRsTestMatch(p, it, "gen"); err == nil {
		t.Error("read-only source: buildLoopRsTestMatch err = nil, want write error")
	}
}

// TestFinalizeRsBacking_RestoreErrorRecorded covers the non-terminal branch where
// restoreRsSource fails (missing backing) and the error is recorded on the
// measurement rather than swallowed.
func TestFinalizeRsBacking_RestoreErrorRecorded(t *testing.T) {
	root := writeRsProj(t, map[string]string{"src/calc.rs": "pub fn add() {}\n"})
	p := funcPayload{Root: root, Fn: model.Function{File: "src/calc.rs", Name: "add"}}
	it := &quest.Item{Key: "src.calc.add"}

	m := &measurement{TestFailed: true} // non-terminal -> attempt restore
	finalizeRsBacking(p, it, m, filepath.Join(".tsma", "test", "orig-absent.rs"))
	if !m.TestFailed || !strings.Contains(m.FailOutput, "failed to restore source") {
		t.Errorf("restore error not recorded: TestFailed=%v FailOutput=%q", m.TestFailed, m.FailOutput)
	}
}

// TestPrepareLoopRs_BuildError covers the buildLoopRsTestMatch-error branch: a
// valid generation but a missing source file makes the inject step fail.
func TestPrepareLoopRs_BuildError(t *testing.T) {
	root := t.TempDir() // no src/calc.rs on disk
	p := funcPayload{Root: root, Fn: model.Function{File: "src/calc.rs", Name: "add", QualifiedName: "src.calc.add"}}
	it := &quest.Item{Key: "src.calc.add"}
	gen := "#[cfg(test)]\nmod tests {\n#[test]\nfn add_works() { assert_eq!(add(1,2),3); }\n}\n"

	m := prepareLoopRs(it, p, []byte(gen))
	if !m.TestFailed || m.Truncated {
		t.Errorf("build error: TestFailed=%v Truncated=%v, want TestFailed only", m.TestFailed, m.Truncated)
	}
	if m.FailOutput == "" {
		t.Error("build error: FailOutput empty, want the read error message")
	}
}
