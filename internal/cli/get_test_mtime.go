//ff:func feature=cli type=helper control=sequence
//ff:what Returns the mtime of a test file in RFC3339 format
package cli

import "time"

// getTestMtime returns the mtime of a single test file in RFC3339 format, or ""
// if the file is missing.
func getTestMtime(root, testFile string) string {
	t := testMtimeOf(root, testFile)
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
