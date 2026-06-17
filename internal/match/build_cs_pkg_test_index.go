//ff:func feature=match type=implementation control=iteration dimension=1 lang=csharp
//ff:what BuildCsPkgTestIndex: parses every *Tests.cs/*Test.cs file in the test directories of srcPkgDir (csTestDirs — the same same-dir + parallel *.Tests SSOT the filename matcher uses) with tree-sitter — via ingestCsTestDir — and builds the name→test-file index. The C# analogue of BuildJavaPkgTestIndex. Returns nil when tree-sitter is unavailable or nothing was indexed, signaling MatchFunc to fall back to filename matching.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// BuildCsPkgTestIndex builds the content-aware test index for srcPkgDir
// (relative to projectRoot) by scanning its candidate test directories. It
// returns nil when the tree-sitter CLI is absent or no test file yielded any
// reference.
func BuildCsPkgTestIndex(projectRoot, srcPkgDir string) *CsPkgTestIndex {
	command := treesitter.ResolveCommand()
	if command == "" {
		return nil
	}
	grammar := treesitter.ResolveGrammar("csharp")

	idx := &CsPkgTestIndex{refs: make(map[string][]string)}
	for _, dir := range csTestDirs(srcPkgDir) {
		ingestCsTestDir(idx, projectRoot, dir, command, grammar)
	}
	if len(idx.refs) == 0 {
		return nil
	}
	return idx
}
