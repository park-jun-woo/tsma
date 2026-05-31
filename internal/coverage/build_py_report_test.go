package coverage

import "testing"

func TestBuildPyReport(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"handler.py": {
				ExecutedLines:    []int{10, 11, 12},
				MissingLines:     []int{13},
				ExecutedBranches: [][]int{{12, 15}},
				MissingBranches:  [][]int{{12, 13}},
			},
		},
	}

	ranges := []pyFuncRange{
		{file: "handler.py", startLine: 10, endLine: 20, funcName: "FuncA"},
	}

	report := buildPyReport(ranges, covData, "/project")

	if len(report.Funcs) != 1 {
		t.Fatalf("got %d funcs, want 1", len(report.Funcs))
	}
	if report.AllCovered {
		t.Error("AllCovered should be false")
	}
	if report.TotalPct == 0 {
		t.Error("TotalPct should be > 0")
	}
}

func TestBuildPyReportAllCovered(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"handler.py": {
				ExecutedLines: []int{10, 11},
			},
		},
	}

	ranges := []pyFuncRange{
		{file: "handler.py", startLine: 10, endLine: 11, funcName: "FuncA"},
	}

	report := buildPyReport(ranges, covData, "/project")

	if !report.AllCovered {
		t.Error("AllCovered should be true")
	}
	if report.TotalPct != 100 {
		t.Errorf("TotalPct = %f, want 100", report.TotalPct)
	}
}

func TestBuildPyReportEmpty(t *testing.T) {
	report := buildPyReport(nil, nil, "/project")

	if !report.AllCovered {
		t.Error("AllCovered should be true for empty ranges")
	}
	if report.TotalPct != 0 {
		t.Errorf("TotalPct = %f, want 0", report.TotalPct)
	}
}
