//ff:func feature=match type=implementation control=sequence lang=java
//ff:what BuildJavaPkgTestIndex: parses every *Test.java/*Tests.java file in the JUnit mirror of srcPkgDir (javaTestDir, the same src/main→src/test SSOT the filename matcher uses) with tree-sitter — via ingestJavaTestDir — and builds the name→test-file index. The Java analogue of BuildTSPkgTestIndex. Returns nil when tree-sitter is unavailable or nothing was indexed, signaling MatchFunc to fall back to filename matching.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// BuildJavaPkgTestIndex builds the content-aware test index for srcPkgDir
// (relative to projectRoot) by scanning its JUnit test mirror directory. It
// returns nil when the tree-sitter CLI is absent or no test file yielded any
// reference.
func BuildJavaPkgTestIndex(projectRoot, srcPkgDir string) *JavaPkgTestIndex {
	command := treesitter.ResolveCommand()
	if command == "" {
		return nil
	}
	grammar := treesitter.ResolveGrammar("java")

	idx := &JavaPkgTestIndex{refs: make(map[string][]string)}
	ingestJavaTestDir(idx, projectRoot, javaTestDir(srcPkgDir), command, grammar)
	if len(idx.refs) == 0 {
		return nil
	}
	return idx
}
