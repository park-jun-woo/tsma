//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Parses one test file and appends its identifier->test references to the index
package match

// ingestTestFile parses the test file at absPath and records, for every
// identifier referenced by each TestXxx in the file, a testRef pointing back to
// that test (using rel as the stored path). Files that fail to parse are
// skipped so one bad file does not abort the whole package index.
func ingestTestFile(idx *PkgTestIndex, absPath, rel string) {
	funcs, err := parseTestFileFuncs(absPath)
	if err != nil {
		return
	}
	for testFunc, idents := range testRefsInFile(funcs) {
		for ident := range idents {
			idx.refs[ident] = append(idx.refs[ident], testRef{File: rel, TestFunc: testFunc})
		}
	}
}
