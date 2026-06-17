//ff:func feature=match type=implementation control=sequence lang=rust
//ff:what BuildRsTestIndex: builds the content-aware test index for one Rust source file (relative to projectRoot) by harvesting (a) the names called inside the file's own in-file #[cfg(test)] modules (attributed to the source file) and (b) the names called by tests/*.rs integration tests (attributed to those files). The Rust analogue of BuildCsPkgTestIndex. Returns nil when the tree-sitter CLI/grammar is unavailable or nothing was indexed, signaling MatchFunc to fall back to filename matching (RsMatcher).
package match

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// BuildRsTestIndex builds the content-aware test index for srcFileRel (relative
// to projectRoot). It returns nil when tree-sitter is absent or no reference was
// found.
func BuildRsTestIndex(projectRoot, srcFileRel string) *RsTestIndex {
	command := treesitter.ResolveCommand()
	if command == "" {
		return nil
	}
	grammar := treesitter.ResolveGrammar("rust")

	idx := &RsTestIndex{refs: make(map[string][]string)}
	ingestRsInFile(idx, command, grammar, filepath.Join(projectRoot, srcFileRel), srcFileRel)
	ingestRsTestsDir(idx, projectRoot, command, grammar)
	if len(idx.refs) == 0 {
		return nil
	}
	return idx
}
