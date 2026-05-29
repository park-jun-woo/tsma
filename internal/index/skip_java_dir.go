//ff:func feature=index type=helper control=selection
//ff:what Returns SkipDir for directories excluded from Java source indexing
package index

import "path/filepath"

// skipJavaDir returns SkipDir for directories that should be excluded from Java
// indexing (build output, VCS, and tsma working dirs).
func skipJavaDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "target", "build", ".gradle", ".git", ".tsma":
		return filepath.SkipDir
	}
	return nil
}
