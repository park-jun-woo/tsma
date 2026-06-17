//ff:func feature=index type=helper control=sequence
//ff:what ResolveCommand: locates the tree-sitter executable for the subprocess pipeline. Honors the TSMA_TREE_SITTER override (absolute path or PATH name); returns "" when absent so every consumer can degrade to its line-based fallback (the tool is an optional prerequisite for indexing).
package treesitter

import (
	"os"
	"os/exec"
	"path/filepath"
)

// ResolveCommand returns the path to a usable tree-sitter executable, or "" if
// none is found. TSMA_TREE_SITTER overrides the default "tree-sitter"; an
// absolute override is used verbatim if it exists, otherwise resolved via PATH.
func ResolveCommand() string {
	name := os.Getenv("TSMA_TREE_SITTER")
	if name == "" {
		name = "tree-sitter"
	}
	if filepath.IsAbs(name) {
		if fi, err := os.Stat(name); err == nil && !fi.IsDir() {
			return name
		}
		return ""
	}
	if p, err := exec.LookPath(name); err == nil {
		return p
	}
	return ""
}
