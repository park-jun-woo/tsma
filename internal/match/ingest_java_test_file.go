//ff:func feature=match type=helper control=iteration dimension=1 lang=java
//ff:what ingestJavaTestFile: parses one JUnit test file with tree-sitter and records, for every called name, a back-reference to that file in the index. The Java analogue of ingestTSTestFile; files that fail to parse are skipped so one bad file never aborts the package index.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// ingestJavaTestFile parses the test file at absPath and appends rel to every
// called name's bucket in idx. Parse failures are silently skipped.
func ingestJavaTestFile(idx *JavaPkgTestIndex, command, grammar, absPath, rel string) {
	root, err := treesitter.ParseFile(command, grammar, absPath)
	if err != nil || root == nil {
		return
	}
	for name := range collectJavaCalledNames(root) {
		idx.refs[name] = append(idx.refs[name], rel)
	}
}
