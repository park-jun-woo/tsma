package coverage

import "testing"

func TestBuildGoReport(t *testing.T) {
	blocks := []coverBlock{
		{file: "github.com/example/pkg/handler.go", startLine: 10, endLine: 15, count: 1},
		{file: "github.com/example/pkg/handler.go", startLine: 16, endLine: 20, count: 0},
		{file: "github.com/example/pkg/handler.go", startLine: 30, endLine: 35, count: 1},
	}

	ranges := []funcRange{
		{file: "pkg/handler.go", startLine: 10, endLine: 20, funcName: "FuncA"},
		{file: "pkg/handler.go", startLine: 30, endLine: 40, funcName: "FuncB"},
	}

	report := buildGoReport(ranges, blocks, "/project")

	if len(report.Funcs) != 2 {
		t.Fatalf("got %d funcs, want 2", len(report.Funcs))
	}
	if report.AllCovered {
		t.Error("AllCovered should be false (FuncA has uncovered block)")
	}
	// FuncA: 2 blocks, 1 covered => 50%; FuncB: 1 block, 1 covered => 100%
	// Total: 3 blocks, 2 covered => 66.66%
	if report.TotalPct < 66.0 || report.TotalPct > 67.0 {
		t.Errorf("TotalPct = %f, want ~66.67", report.TotalPct)
	}
	if len(report.Uncovered) != 1 {
		t.Errorf("Uncovered = %d, want 1", len(report.Uncovered))
	}
}

func TestBuildGoReportAllCovered(t *testing.T) {
	blocks := []coverBlock{
		{file: "github.com/example/pkg/handler.go", startLine: 10, endLine: 15, count: 1},
	}
	ranges := []funcRange{
		{file: "pkg/handler.go", startLine: 10, endLine: 15, funcName: "Covered"},
	}

	report := buildGoReport(ranges, blocks, "/project")

	if !report.AllCovered {
		t.Error("AllCovered should be true")
	}
	if report.TotalPct != 100 {
		t.Errorf("TotalPct = %f, want 100", report.TotalPct)
	}
}

func TestBuildGoReportEmpty(t *testing.T) {
	report := buildGoReport(nil, nil, "/project")

	if !report.AllCovered {
		t.Error("AllCovered should be true for empty ranges")
	}
	if report.TotalPct != 0 {
		t.Errorf("TotalPct = %f, want 0 (no blocks)", report.TotalPct)
	}
}
