//ff:func feature=match type=helper control=iteration dimension=1 lang=rust
//ff:what rsRefsToTestMatch: collapses the test-file list for a function into a TestMatch with deduplicated, order-preserving Files. TestFuncs is left nil because cargo runs whole test binaries (the runner resolves cases), exactly as the non-Go fallback convention specifies. Mirrors csRefsToTestMatch.
package match

// rsRefsToTestMatch deduplicates files (order-preserving) into a TestMatch. It
// reports found=false when no files result.
func rsRefsToTestMatch(files []string) (TestMatch, bool) {
	seen := make(map[string]struct{})
	var out []string
	for _, f := range files {
		if _, dup := seen[f]; dup {
			continue
		}
		seen[f] = struct{}{}
		out = append(out, f)
	}
	if len(out) == 0 {
		return TestMatch{}, false
	}
	return TestMatch{Files: out, TestFuncs: nil}, true
}
