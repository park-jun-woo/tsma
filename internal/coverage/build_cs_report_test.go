package coverage

import "testing"

func TestBuildCsReportFullyCovered(t *testing.T) {
	cov, err := parseCobertura(coberturaFixture)
	if err != nil {
		t.Fatal(err)
	}
	ranges := []csFuncRange{{file: "App/Calculator.cs", startLine: 5, endLine: 5, funcName: "Add"}}
	report := buildCsReport(ranges, cov, "/proj")

	if !report.AllCovered {
		t.Error("expected AllCovered true for fully covered function")
	}
	if report.TotalPct != 100 {
		t.Errorf("TotalPct = %v, want 100", report.TotalPct)
	}
	if len(report.Uncovered) != 0 {
		t.Errorf("Uncovered = %v, want none", report.Uncovered)
	}
}

func TestBuildCsReportPartial(t *testing.T) {
	cov, err := parseCobertura(coberturaFixture)
	if err != nil {
		t.Fatal(err)
	}
	ranges := []csFuncRange{{file: "App/Calculator.cs", startLine: 9, endLine: 12, funcName: "Classify"}}
	report := buildCsReport(ranges, cov, "/proj")

	if report.AllCovered {
		t.Error("expected AllCovered false for partially covered function")
	}
	if report.TotalPct <= 0 || report.TotalPct >= 100 {
		t.Errorf("TotalPct = %v, want strictly between 0 and 100", report.TotalPct)
	}
	if len(report.Uncovered) == 0 {
		t.Error("expected uncovered branches/lines to be recorded")
	}
}

func TestBuildCsReportEmpty(t *testing.T) {
	report := buildCsReport(nil, nil, "/proj")
	if !report.AllCovered {
		t.Error("AllCovered should be true for empty ranges")
	}
	if report.TotalPct != 0 {
		t.Errorf("TotalPct = %v, want 0", report.TotalPct)
	}
}
