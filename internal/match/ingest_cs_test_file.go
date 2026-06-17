//ff:func feature=match type=helper control=iteration dimension=1 lang=csharp
//ff:what ingestCsTestFile: parses one C# test file with tree-sitter and records, for every called name, a back-reference to that file in the index. The C# analogue of ingestJavaTestFile; files that fail to parse are skipped so one bad file never aborts the index.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// ingestCsTestFile parses the test file at absPath and appends rel to every
// called name's bucket in idx. Parse failures are silently skipped.
func ingestCsTestFile(idx *CsPkgTestIndex, command, grammar, absPath, rel string) {
	root, err := treesitter.ParseFile(command, grammar, absPath)
	if err != nil || root == nil {
		return
	}
	for name := range collectCsCalledNames(root) {
		idx.refs[name] = append(idx.refs[name], rel)
	}
}
