package match

import (
	"os"
	"path/filepath"
	"testing"
)

// writeFakeTreeSitter writes a stub tree-sitter CLI that ignores its input file
// and prints a canned --xml parse tree (body) wrapped in <sources><source
// name="<last arg>">, mirroring the real CLI's output shape. It returns the
// script's absolute path, suitable for the TSMA_TREE_SITTER override, so the
// content-index branches can be exercised deterministically without the real
// CLI or grammars (same stub-executable convention as the fake Python
// interpreter in py_index_helpers_test.go).
func writeFakeTreeSitter(t *testing.T, body string) string {
	t.Helper()
	script := filepath.Join(t.TempDir(), "fake-tree-sitter")
	src := "#!/bin/sh\n" +
		"for a in \"$@\"; do last=\"$a\"; done\n" +
		"printf '<sources><source name=\"%s\">' \"$last\"\n" +
		"cat <<'XMLEOF'\n" +
		body + "\n" +
		"XMLEOF\n" +
		"printf '</source></sources>'\n"
	if err := os.WriteFile(script, []byte(src), 0o755); err != nil {
		t.Fatalf("write fake tree-sitter: %v", err)
	}
	return script
}

// useFakeTreeSitter points the resolver env at the stub CLI and pins the
// language grammar override to an existing directory so ResolveGrammar never
// probes the host's node_modules (deterministic regardless of what is
// installed).
func useFakeTreeSitter(t *testing.T, grammarEnv, body string) {
	t.Helper()
	t.Setenv("TSMA_TREE_SITTER", writeFakeTreeSitter(t, body))
	t.Setenv(grammarEnv, t.TempDir())
}
