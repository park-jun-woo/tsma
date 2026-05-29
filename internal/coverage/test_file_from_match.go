//ff:func feature=coverage type=helper control=sequence
//ff:what Returns the representative single test file from a TestMatch, or "" if none
package coverage

import "github.com/park-jun-woo/tsma/internal/match"

// testFileFromMatch returns the representative test file for non-Go checkers,
// which match a single file per function. It is the first file in m.Files, or
// "" when the match carries no files.
func testFileFromMatch(m match.TestMatch) string {
	if len(m.Files) == 0 {
		return ""
	}
	return m.Files[0]
}
