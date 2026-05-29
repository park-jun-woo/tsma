//ff:func feature=match type=helper control=iteration dimension=1 lang=go
//ff:what Deduplicates test refs into a TestMatch of distinct files and test funcs
package match

// refsToTestMatch collapses a slice of testRefs into a TestMatch with
// deduplicated, order-preserving Files and TestFuncs. It reports found=false
// when no files result (which is treated as "no attribution").
func refsToTestMatch(refs []testRef) (TestMatch, bool) {
	seenFiles := make(map[string]struct{})
	seenFuncs := make(map[string]struct{})
	var files, funcs []string
	for _, r := range refs {
		if _, dup := seenFiles[r.File]; !dup {
			seenFiles[r.File] = struct{}{}
			files = append(files, r.File)
		}
		if _, dup := seenFuncs[r.TestFunc]; !dup {
			seenFuncs[r.TestFunc] = struct{}{}
			funcs = append(funcs, r.TestFunc)
		}
	}
	if len(files) == 0 {
		return TestMatch{}, false
	}
	return TestMatch{Files: files, TestFuncs: funcs}, true
}
