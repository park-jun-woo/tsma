//ff:func feature=cli type=helper control=sequence
//ff:what Guesses the expected test file path from a source file path
package cli

import (
	"path/filepath"
	"strings"
)

// guessTestFile guesses the expected test file path from a source file path.
func guessTestFile(srcFile string) string {
	dir := filepath.Dir(srcFile)
	base := filepath.Base(srcFile)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	return filepath.Join(dir, name+"_test"+ext)
}
