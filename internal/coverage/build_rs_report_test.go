package coverage

import "testing"

func TestBuildRsReportFullyCovered(t *testing.T) {
	cov, err := parseLLVMCov(llvmCovFixture)
	if err != nil {
		t.Fatal(err)
	}
	ranges := []rsFuncRange{{file: "src/lib.rs", startLine: 1, endLine: 5, funcName: "classify"}}
	report := buildRsReport(ranges, cov, "/home/dev/demo")

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

func TestBuildRsReportPartial(t *testing.T) {
	cov, err := parseLLVMCov(llvmCovFixture)
	if err != nil {
		t.Fatal(err)
	}
	ranges := []rsFuncRange{{file: "src/lib.rs", startLine: 8, endLine: 10, funcName: "unused"}}
	report := buildRsReport(ranges, cov, "/home/dev/demo")

	if report.AllCovered {
		t.Error("expected AllCovered false for partially covered function")
	}
	if report.TotalPct != 50 {
		t.Errorf("TotalPct = %v, want 50", report.TotalPct)
	}
	if len(report.Uncovered) == 0 {
		t.Error("expected uncovered branches to be recorded")
	}
}
