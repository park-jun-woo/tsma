//ff:func feature=detect type=helper control=sequence
//ff:what Checks if a file exists at the given path
package detect

import "os"

// fileExists reports whether a file exists at path. This mirrors the
// runner-package helper of the same name; detect owns the pytest-detection SSOT
// so the helper is duplicated here to avoid an import cycle (detect must not
// depend on runner/coverage).
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
