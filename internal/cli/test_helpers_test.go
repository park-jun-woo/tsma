package cli

import "github.com/park-jun-woo/tsma/internal/match"

// mkMatch wraps a single test file into a TestMatch with nil TestFuncs, the
// shape produced by the non-Go fallback matcher and sufficient for the Go path
// (which resolves test functions from the file when TestFuncs is nil). Used by
// cli tests to drive the new TestMatch-based runAndMeasure signature.
func mkMatch(file string) match.TestMatch {
	return match.TestMatch{Files: []string{file}}
}
