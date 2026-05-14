package coverage

import "testing"

func TestBuildTSReport(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"src/handler.ts": {
			StatementMap: map[string]coverageRange{
				"0": {Start: coveragePosition{Line: 10}, End: coveragePosition{Line: 10}},
				"1": {Start: coveragePosition{Line: 15}, End: coveragePosition{Line: 15}},
			},
			S:         map[string]int{"0": 1, "1": 0},
			BranchMap: map[string]coverageBranch{},
			B:         map[string][]int{},
			FnMap:     map[string]coverageFunction{},
			F:         map[string]int{},
		},
	}

	ranges := []tsFuncRange{
		{file: "src/handler.ts", startLine: 10, endLine: 20, funcName: "handler"},
	}

	report := buildTSReport(ranges, data, "/project")

	if len(report.Funcs) != 1 {
		t.Fatalf("got %d funcs, want 1", len(report.Funcs))
	}
	if report.AllCovered {
		t.Error("AllCovered should be false")
	}
	if report.TotalPct != 50 {
		t.Errorf("TotalPct = %f, want 50", report.TotalPct)
	}
}

func TestBuildTSReportAllCovered(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"src/handler.ts": {
			StatementMap: map[string]coverageRange{
				"0": {Start: coveragePosition{Line: 10}, End: coveragePosition{Line: 10}},
			},
			S:         map[string]int{"0": 1},
			BranchMap: map[string]coverageBranch{},
			B:         map[string][]int{},
			FnMap:     map[string]coverageFunction{},
			F:         map[string]int{},
		},
	}

	ranges := []tsFuncRange{
		{file: "src/handler.ts", startLine: 10, endLine: 20, funcName: "handler"},
	}

	report := buildTSReport(ranges, data, "/project")

	if !report.AllCovered {
		t.Error("AllCovered should be true")
	}
}

func TestBuildTSReportEmpty(t *testing.T) {
	report := buildTSReport(nil, nil, "/project")

	if !report.AllCovered {
		t.Error("AllCovered should be true for empty ranges")
	}
	if report.TotalPct != 0 {
		t.Errorf("TotalPct = %f, want 0", report.TotalPct)
	}
}
