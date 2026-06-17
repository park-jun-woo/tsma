//ff:func feature=index type=helper control=iteration dimension=1
//ff:what ResolveGrammar: locates the tree-sitter grammar directory for a language, passed to `tree-sitter parse -p`. Honors a per-language TSMA_<LANG>_GRAMMAR override and otherwise probes common node_modules install bases. The one place 005b~e register their grammar (add a langGrammar entry); returns "" to let the CLI resolve via configured parser dirs (and ultimately fall back to line-based).
package treesitter

import (
	"os"
	"path/filepath"
)

// langGrammar maps a tsma language label to its grammar override env var and the
// node_modules package path (relative to a base dir) that holds the grammar.
var langGrammar = map[string]struct {
	env     string
	nodePkg []string
}{
	"typescript": {env: "TSMA_TS_GRAMMAR", nodePkg: []string{"tree-sitter-typescript", "typescript"}},
	"java":       {env: "TSMA_JAVA_GRAMMAR", nodePkg: []string{"tree-sitter-java"}},
	"csharp":     {env: "TSMA_CSHARP_GRAMMAR", nodePkg: []string{"tree-sitter-c-sharp"}},
	"rust":       {env: "TSMA_RUST_GRAMMAR", nodePkg: []string{"tree-sitter-rust"}},
}

// ResolveGrammar returns the grammar directory for lang, or "" if not found.
// The per-language env override wins; otherwise common install bases (cwd, /tmp,
// home) are probed under node_modules.
func ResolveGrammar(lang string) string {
	cfg, ok := langGrammar[lang]
	if !ok {
		return ""
	}
	if g := os.Getenv(cfg.env); g != "" {
		return envGrammarDir(g)
	}
	for _, base := range grammarCandidateBases() {
		parts := append([]string{base, "node_modules"}, cfg.nodePkg...)
		dir := filepath.Join(parts...)
		if grammarDirExists(dir) {
			return dir
		}
	}
	return ""
}
