//ff:func feature=index type=helper control=sequence
//ff:what grammarCandidateBases: lists directories whose node_modules may hold a tree-sitter grammar (cwd, /tmp, home) — probed when no env override is set.
package treesitter

import "os"

// grammarCandidateBases lists directories whose node_modules may hold a grammar.
func grammarCandidateBases() []string {
	bases := []string{".", "/tmp"}
	if home, err := os.UserHomeDir(); err == nil {
		bases = append(bases, home)
	}
	return bases
}
