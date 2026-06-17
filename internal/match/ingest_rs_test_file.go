//ff:func feature=match type=helper control=iteration dimension=1 lang=rust
//ff:what ingestRsTestFile: parses one tests/ integration file with tree-sitter and records, for every called name (whole-file test code via collectRsCalledNames), a back-reference to that file in the index. The Rust analogue of ingestCsTestFile; files that fail to parse are skipped so one bad file never aborts the index.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// ingestRsTestFile parses the integration test at absPath and appends rel to
// every called name's bucket in idx. Parse failures are silently skipped.
func ingestRsTestFile(idx *RsTestIndex, command, grammar, absPath, rel string) {
	root, err := treesitter.ParseFile(command, grammar, absPath)
	if err != nil || root == nil {
		return
	}
	for name := range collectRsCalledNames(root) {
		idx.refs[name] = append(idx.refs[name], rel)
	}
}
