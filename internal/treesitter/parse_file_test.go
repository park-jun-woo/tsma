package treesitter

import (
	"os"
	"path/filepath"
	"testing"
)

// fakeTreeSitter writes an executable shell script into a temp dir and returns
// its path, so ParseFile can be driven without a real tree-sitter binary.
func fakeTreeSitter(t *testing.T, script string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "fake-tree-sitter")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+script+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseFileMatchByName(t *testing.T) {
	// the script echoes XML whose source name is the last argument (the file
	// path Run appends), so the name-match branch returns that root.
	cmd := fakeTreeSitter(t, `for last; do :; done
printf '<sources><source name="%s"><program srow="0" scol="0" erow="1" ecol="0"><identifier srow="0" scol="0" erow="0" ecol="1">x</identifier></program></source></sources>' "$last"`)
	root, err := ParseFile(cmd, "", "/abs/x.ts")
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if root == nil || root.Type != "program" {
		t.Fatalf("root = %+v, want program", root)
	}
}

func TestParseFileSingleSourceFallback(t *testing.T) {
	// name does not match the asked path, but a sole source with a root is
	// still returned via the single-source fallback.
	cmd := fakeTreeSitter(t, `printf '<sources><source name="/other.ts"><program srow="0" scol="0" erow="1" ecol="0"></program></source></sources>'`)
	root, err := ParseFile(cmd, "", "/abs/x.ts")
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if root == nil || root.Type != "program" {
		t.Fatalf("root = %+v, want program", root)
	}
}

func TestParseFileMalformedXML(t *testing.T) {
	cmd := fakeTreeSitter(t, `printf '<sources><a></b></sources>'`)
	if _, err := ParseFile(cmd, "", "/abs/x.ts"); err == nil {
		t.Error("expected ParseXML error for malformed XML")
	}
}

func TestParseFileNoParseTree(t *testing.T) {
	cmd := fakeTreeSitter(t, `printf '<sources></sources>'`)
	if _, err := ParseFile(cmd, "", "/abs/x.ts"); err == nil {
		t.Error("expected no-parse-tree error for empty sources")
	}
}
