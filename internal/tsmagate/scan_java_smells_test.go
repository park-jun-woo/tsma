package tsmagate

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// locateJavaTSGate finds a usable tree-sitter CLI + tree-sitter-java grammar and
// exports the env vars smell.ScanJava reads. Returns false -> caller t.Skip.
func locateJavaTSGate(t *testing.T) bool {
	t.Helper()
	bases := []string{".", "..", "../..", "/tmp"}
	if home, err := os.UserHomeDir(); err == nil {
		bases = append(bases, home)
	}
	cmd := os.Getenv("TSMA_TREE_SITTER")
	if cmd == "" {
		if p, err := exec.LookPath("tree-sitter"); err == nil {
			cmd = p
		}
	}
	if cmd == "" {
		for _, b := range bases {
			p := filepath.Join(b, "node_modules", ".bin", "tree-sitter")
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				cmd, _ = filepath.Abs(p)
				break
			}
		}
	}
	grammar := os.Getenv("TSMA_JAVA_GRAMMAR")
	if grammar == "" {
		for _, b := range bases {
			p := filepath.Join(b, "node_modules", "tree-sitter-java")
			if fi, err := os.Stat(p); err == nil && fi.IsDir() {
				grammar, _ = filepath.Abs(p)
				break
			}
		}
	}
	if cmd == "" || grammar == "" {
		return false
	}
	t.Setenv("TSMA_TREE_SITTER", cmd)
	t.Setenv("TSMA_JAVA_GRAMMAR", grammar)
	return true
}

// TestScanJavaSmellsContinue covers the error branch: when tree-sitter is
// unavailable every ScanJava errors and is skipped, yielding no findings.
func TestScanJavaSmellsContinue(t *testing.T) {
	t.Setenv("TSMA_TREE_SITTER", "/nonexistent/abs/tree-sitter")
	got := scanJavaSmells("../../testdata/java", []string{"smell/ReflectionTest.java"})
	if len(got) != 0 {
		t.Errorf("expected no findings when tree-sitter unavailable, got %+v", got)
	}
}

// TestScanJavaSmellsAppend covers the append branch with a real CLI: the
// ReflectionTest fixture yields escape-hatch findings, while a nonexistent file
// errors and is skipped (continue) within the same loop.
func TestScanJavaSmellsAppend(t *testing.T) {
	if !locateJavaTSGate(t) {
		t.Skip("tree-sitter CLI + java grammar not available")
	}
	got := scanJavaSmells("../../testdata/java", []string{
		"smell/ReflectionTest.java",
		"smell/does_not_exist.java", // errors -> continue
	})
	if len(got) == 0 {
		t.Fatal("expected findings from ReflectionTest.java")
	}
	for _, f := range got {
		if f.Rule != "TS-REFL-JV-001" && f.Rule != "TS-REFL-JV-002" {
			t.Errorf("unexpected rule %q", f.Rule)
		}
	}
}
