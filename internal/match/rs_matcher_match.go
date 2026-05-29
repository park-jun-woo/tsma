//ff:func feature=match type=implementation control=sequence
//ff:what Returns the source file itself for in-file #[cfg(test)], else looks for an integration test in tests/
package match

import (
	"os"
	"path/filepath"
	"strings"
)

// Match finds the test file covering the given Rust source file.
//
// Rust unit tests commonly live in the same file inside a `#[cfg(test)]`
// module; in that case the source file is returned as its own test file.
// Otherwise an integration test at `tests/<name>.rs` (relative to the project
// root) is searched. The returned testFile is a projectRoot-relative path.
func (m *RsMatcher) Match(projectRoot string, sourceFile string) (string, bool) {
	absSrc := filepath.Join(projectRoot, sourceFile)

	// In-file unit tests: the source file is its own test file.
	if hasInFileTests(absSrc) {
		return sourceFile, true
	}

	// Integration tests: tests/<name>.rs at the project root.
	name := strings.TrimSuffix(filepath.Base(sourceFile), ".rs")
	testRel := filepath.Join("tests", name+".rs")
	if _, err := os.Stat(filepath.Join(projectRoot, testRel)); err == nil {
		return testRel, true
	}

	return "", false
}
