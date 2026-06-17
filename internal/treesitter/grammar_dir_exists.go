//ff:func feature=index type=helper control=sequence
//ff:what grammarDirExists: reports whether a path is an existing directory (grammar-dir probe predicate).
package treesitter

import "os"

// grammarDirExists reports whether path is an existing directory.
func grammarDirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}
