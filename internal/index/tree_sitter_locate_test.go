package index

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// locateTreeSitter finds a usable tree-sitter CLI + TypeScript grammar for the
// integration tests, probing env, PATH, and common node_modules install bases.
// It returns ok=false when neither is found so callers can t.Skip — the precise
// path is an optional prerequisite, exactly as the plan's skip-gate strategy
// requires. On success it exports TSMA_TREE_SITTER/TSMA_TS_GRAMMAR for the test.
func locateTreeSitter(t *testing.T) (cmd, grammar string, ok bool) {
	t.Helper()

	cmd = os.Getenv("TSMA_TREE_SITTER")
	if cmd == "" {
		if p, err := exec.LookPath("tree-sitter"); err == nil {
			cmd = p
		}
	}
	if cmd == "" {
		cmd = probe(tsLocateBases(), filepath.Join("node_modules", ".bin", "tree-sitter"))
	}

	grammar = os.Getenv("TSMA_TS_GRAMMAR")
	if grammar == "" {
		grammar = probeDir(tsLocateBases(), filepath.Join("node_modules", "tree-sitter-typescript", "typescript"))
	}

	if cmd == "" || grammar == "" {
		return "", "", false
	}
	t.Setenv("TSMA_TREE_SITTER", cmd)
	t.Setenv("TSMA_TS_GRAMMAR", grammar)
	return cmd, grammar, true
}

// tsLocateBases lists directories whose node_modules may hold the CLI/grammar.
func tsLocateBases() []string {
	bases := []string{".", "..", "../..", "/tmp"}
	if home, err := os.UserHomeDir(); err == nil {
		bases = append(bases, home)
	}
	return bases
}

// probe returns the first base/rel that exists as a regular file, or "".
func probe(bases []string, rel string) string {
	for _, b := range bases {
		p := filepath.Join(b, rel)
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			abs, _ := filepath.Abs(p)
			return abs
		}
	}
	return ""
}

// probeDir returns the first base/rel that exists as a directory, or "".
func probeDir(bases []string, rel string) string {
	for _, b := range bases {
		p := filepath.Join(b, rel)
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			abs, _ := filepath.Abs(p)
			return abs
		}
	}
	return ""
}
