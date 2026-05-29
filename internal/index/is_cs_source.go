//ff:func feature=index type=helper control=iteration dimension=1 lang=csharp
//ff:what Returns true if the path is a non-test C# source file eligible for indexing
package index

import "strings"

// isCsSource returns true if the path is a C# source file eligible for
// indexing. Files in a *.Tests project directory and *Test.cs / *Tests.cs
// files are excluded, since indexing targets production code.
func isCsSource(path string) bool {
	if !strings.HasSuffix(path, ".cs") {
		return false
	}
	slashed := strings.ReplaceAll(path, "\\", "/")
	for _, seg := range strings.Split(slashed, "/") {
		if strings.HasSuffix(seg, ".Tests") || strings.HasSuffix(seg, ".Test") {
			return false
		}
	}
	base := slashed
	if idx := strings.LastIndex(slashed, "/"); idx >= 0 {
		base = slashed[idx+1:]
	}
	name := strings.TrimSuffix(base, ".cs")
	if strings.HasSuffix(name, "Test") || strings.HasSuffix(name, "Tests") {
		return false
	}
	return true
}
