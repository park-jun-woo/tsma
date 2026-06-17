//ff:func feature=match type=helper control=iteration dimension=1 lang=csharp
//ff:what csRefsToTestMatch: collapses the test-file list for a function into a TestMatch with deduplicated, order-preserving Files. TestFuncs is left nil because dotnet runs whole test classes (the runner resolves cases), exactly as the non-Go fallback convention specifies. Mirrors javaRefsToTestMatch.
package match

// csRefsToTestMatch deduplicates files (order-preserving) into a TestMatch. It
// reports found=false when no files result.
func csRefsToTestMatch(files []string) (TestMatch, bool) {
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
