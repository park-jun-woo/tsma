package cli

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/match"
)

// TestRepresentativeTestFile_firstFile verifies the first matched file is
// returned as the representative.
func TestRepresentativeTestFile_firstFile(t *testing.T) {
	tm := match.TestMatch{Files: []string{"first_test.go", "second_test.go"}}
	if got := representativeTestFile(tm); got != "first_test.go" {
		t.Errorf("representativeTestFile = %q, want %q", got, "first_test.go")
	}
}

// TestRepresentativeTestFile_emptyMatch verifies an empty match yields "".
func TestRepresentativeTestFile_emptyMatch(t *testing.T) {
	if got := representativeTestFile(match.TestMatch{}); got != "" {
		t.Errorf("representativeTestFile = %q, want empty", got)
	}
}
