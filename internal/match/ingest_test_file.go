//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Parses one test file and appends its identifier->test references to the index
package match

// ingestTestFile parses the test file at absPath and records, for every
// identifier referenced by each TestXxx in the file, a testRef pointing back to
// that test (using rel as the stored path). The index key stays the bare
// identifier name; the call's statically-resolved receiver type is carried on
// the testRef so receiver-aware lookup can filter the bucket. A name called on
// several receivers within one test produces one testRef per (name, receiver)
// pair. Files that fail to parse are skipped so one bad file does not abort the
// whole package index.
func ingestTestFile(idx *PkgTestIndex, absPath, rel string) {
	funcs, err := parseTestFileFuncs(absPath)
	if err != nil {
		return
	}
	for testFunc, refs := range testRefsInFile(funcs) {
		for ref := range refs {
			idx.refs[ref.Name] = append(idx.refs[ref.Name], testRef{
				File:     rel,
				TestFunc: testFunc,
				Receiver: ref.Receiver,
			})
		}
	}
}
