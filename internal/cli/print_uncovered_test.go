package cli

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/coverage"
)

func TestPrintUncovered_multipleBranches(t *testing.T) {
	branches := []coverage.UncoveredBranch{
		{File: "pkg/foo.go", Line: 10},
		{File: "pkg/bar.go", Line: 25},
	}
	output := captureStdout(func() {
		printUncovered(branches)
	})

	if !strings.Contains(output, "UNCOVERED: pkg/foo.go:10") {
		t.Errorf("expected pkg/foo.go:10 in output, got %q", output)
	}
	if !strings.Contains(output, "UNCOVERED: pkg/bar.go:25") {
		t.Errorf("expected pkg/bar.go:25 in output, got %q", output)
	}
}

func TestPrintUncovered_empty(t *testing.T) {
	output := captureStdout(func() {
		printUncovered(nil)
	})
	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}
