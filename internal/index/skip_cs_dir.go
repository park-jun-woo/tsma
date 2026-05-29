//ff:func feature=index type=helper control=selection lang=csharp
//ff:what Returns SkipDir for directories excluded from C# source indexing
package index

import "path/filepath"

// skipCsDir returns SkipDir for directories that should be excluded from C#
// indexing (build output, VCS, and tsma working dirs).
func skipCsDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "bin", "obj", ".vs", ".git", ".tsma":
		return filepath.SkipDir
	}
	return nil
}
