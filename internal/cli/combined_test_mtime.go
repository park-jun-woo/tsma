//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Computes the latest (max) mtime across a set of test files as the change signature
package cli

import "time"

// combinedTestMtime returns the latest mtime across all given test files,
// formatted as RFC3339, to use as the change-detection signature for a function
// covered by multiple test files. If any one file is modified, the max advances
// and the function is re-measured. Missing files contribute nothing; an empty or
// all-missing set yields "".
func combinedTestMtime(root string, testFiles []string) string {
	var latest time.Time
	for _, tf := range testFiles {
		t := testMtimeOf(root, tf)
		if t.After(latest) {
			latest = t
		}
	}
	if latest.IsZero() {
		return ""
	}
	return latest.Format(time.RFC3339)
}
