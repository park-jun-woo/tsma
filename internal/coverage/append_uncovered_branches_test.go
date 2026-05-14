package coverage

import "testing"

func TestAppendUncoveredBranches(t *testing.T) {
	report := &Report{}
	fc := FuncCoverage{
		File:           "handler.go",
		UncoveredLines: []int{10, 15, 20},
	}

	appendUncoveredBranches(report, fc)

	if len(report.Uncovered) != 3 {
		t.Fatalf("got %d uncovered, want 3", len(report.Uncovered))
	}
	for i, line := range []int{10, 15, 20} {
		if report.Uncovered[i].File != "handler.go" {
			t.Errorf("Uncovered[%d].File = %q, want %q", i, report.Uncovered[i].File, "handler.go")
		}
		if report.Uncovered[i].Line != line {
			t.Errorf("Uncovered[%d].Line = %d, want %d", i, report.Uncovered[i].Line, line)
		}
	}
}

func TestAppendUncoveredBranchesEmpty(t *testing.T) {
	report := &Report{}
	fc := FuncCoverage{
		File:           "handler.go",
		UncoveredLines: nil,
	}

	appendUncoveredBranches(report, fc)

	if len(report.Uncovered) != 0 {
		t.Errorf("expected no uncovered entries, got %d", len(report.Uncovered))
	}
}

func TestAppendUncoveredBranchesAccumulate(t *testing.T) {
	report := &Report{
		Uncovered: []UncoveredBranch{{File: "old.go", Line: 1}},
	}
	fc := FuncCoverage{
		File:           "new.go",
		UncoveredLines: []int{5},
	}

	appendUncoveredBranches(report, fc)

	if len(report.Uncovered) != 2 {
		t.Fatalf("got %d uncovered, want 2", len(report.Uncovered))
	}
	if report.Uncovered[0].File != "old.go" {
		t.Errorf("first entry should be old.go")
	}
	if report.Uncovered[1].File != "new.go" || report.Uncovered[1].Line != 5 {
		t.Errorf("second entry should be new.go:5, got %s:%d", report.Uncovered[1].File, report.Uncovered[1].Line)
	}
}
