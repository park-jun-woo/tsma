package runner

import "github.com/park-jun-woo/tsma/internal/match"

// mkMatch wraps a single test file into a TestMatch with nil TestFuncs, the
// shape produced by the non-Go fallback matcher. Used by runner tests to drive
// the new TestMatch-based Run signature.
func mkMatch(file string) match.TestMatch {
	return match.TestMatch{Files: []string{file}}
}
