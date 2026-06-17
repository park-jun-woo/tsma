//ff:func feature=match type=implementation control=iteration dimension=1 lang=typescript
//ff:what BuildTSPkgTestIndex: parses every *.test/*.spec file in the package dir (and its __tests__/) with tree-sitter — via ingestTSDir per directory — and builds the name→test-file index. The TS analogue of BuildPkgTestIndex. Returns nil when tree-sitter is unavailable or nothing was indexed, signaling MatchFunc to fall back to filename matching.
package match

import (
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/treesitter"
)

// BuildTSPkgTestIndex builds the content-aware test index for pkgDir (relative
// to projectRoot), scanning pkgDir and pkgDir/__tests__. It returns nil when the
// tree-sitter CLI is absent or no test file yielded any reference.
func BuildTSPkgTestIndex(projectRoot, pkgDir string) *TSPkgTestIndex {
	command := treesitter.ResolveCommand()
	if command == "" {
		return nil
	}
	grammar := treesitter.ResolveGrammar("typescript")

	idx := &TSPkgTestIndex{refs: make(map[string][]string)}
	for _, dir := range []string{pkgDir, filepath.Join(pkgDir, "__tests__")} {
		ingestTSDir(idx, projectRoot, dir, command, grammar)
	}
	if len(idx.refs) == 0 {
		return nil
	}
	return idx
}
