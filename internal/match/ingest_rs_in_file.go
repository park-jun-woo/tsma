//ff:func feature=match type=helper control=iteration dimension=1 lang=rust
//ff:what ingestRsInFile: parses the source file and records, for every name called inside its in-file #[cfg(test)] modules (rsCfgTestBodies → collectRsCalledNames), a back-reference to the source file itself in the index — the in-file unit-test attribution (the source file IS its own test file in Rust). A parse failure is a no-op so a single bad file never aborts the index.
package match

import "github.com/park-jun-woo/tsma/internal/treesitter"

// ingestRsInFile adds the source file as a test reference for every function its
// own #[cfg(test)] module calls.
func ingestRsInFile(idx *RsTestIndex, command, grammar, absPath, srcRel string) {
	root, err := treesitter.ParseFile(command, grammar, absPath)
	if err != nil || root == nil {
		return
	}
	for _, body := range rsCfgTestBodies(root) {
		for name := range collectRsCalledNames(body) {
			idx.refs[name] = append(idx.refs[name], srcRel)
		}
	}
}
