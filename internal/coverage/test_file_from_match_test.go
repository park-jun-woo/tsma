package coverage

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/match"
)

func TestTestFileFromMatch_first(t *testing.T) {
	got := testFileFromMatch(match.TestMatch{Files: []string{"a_test.go", "b_test.go"}})
	if got != "a_test.go" {
		t.Errorf("got %q, want a_test.go", got)
	}
}

func TestTestFileFromMatch_empty(t *testing.T) {
	if got := testFileFromMatch(match.TestMatch{}); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}
