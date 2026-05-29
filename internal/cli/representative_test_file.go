//ff:func feature=cli type=helper control=sequence
//ff:what Returns the representative (first) test file of a match for display, or "" if none
package cli

import "github.com/park-jun-woo/tsma/internal/match"

// representativeTestFile returns the first matched test file, used as the
// display/back-compat representative value. It is "" when the match is empty.
func representativeTestFile(tm match.TestMatch) string {
	if len(tm.Files) == 0 {
		return ""
	}
	return tm.Files[0]
}
