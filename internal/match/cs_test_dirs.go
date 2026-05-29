//ff:func feature=match type=helper control=sequence lang=csharp
//ff:what Returns candidate test directories for a C# source directory (same dir and *.Tests project)
package match

import "strings"

// csTestDirs returns the directories where a C# test file for a source under
// srcDir might live, in priority order:
//
//  1. the same directory (tests beside sources), and
//  2. the parallel *.Tests project directory, formed by suffixing the first
//     path segment (the project directory) with ".Tests". For example
//     "App/Services" yields "App.Tests/Services".
//
// All returned paths use the OS separator-agnostic slash form on input but are
// returned with the original separators preserved by the caller's filepath.Join.
func csTestDirs(srcDir string) []string {
	dirs := []string{srcDir}

	slashed := strings.ReplaceAll(srcDir, "\\", "/")
	if slashed == "" || slashed == "." {
		return dirs
	}

	parts := strings.SplitN(slashed, "/", 2)
	proj := parts[0] + ".Tests"
	if len(parts) == 2 {
		dirs = append(dirs, proj+"/"+parts[1])
	} else {
		dirs = append(dirs, proj)
	}
	return dirs
}
