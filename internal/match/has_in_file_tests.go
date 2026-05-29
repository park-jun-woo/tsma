//ff:func feature=match type=helper control=sequence
//ff:what Returns true if a Rust source file contains an in-file #[cfg(test)] module
package match

import (
	"os"
	"strings"
)

// hasInFileTests reports whether the given Rust source file at absPath contains
// an inline `#[cfg(test)]` attribute, indicating in-file unit tests.
func hasInFileTests(absPath string) bool {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "#[cfg(test)]")
}
