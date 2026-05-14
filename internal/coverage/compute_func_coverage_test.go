package coverage

import "testing"

func TestComputeFuncCoveragePartial(t *testing.T) {
	blocks := []coverBlock{
		{file: "github.com/example/pkg/handler.go", startLine: 10, endLine: 15, count: 1},
		{file: "github.com/example/pkg/handler.go", startLine: 16, endLine: 20, count: 0},
	}

	r := funcRange{file: "pkg/handler.go", startLine: 10, endLine: 25, funcName: "Handler"}
	fc := computeFuncCoverage(r, blocks, "/project")

	if fc.TotalBlocks != 2 {
		t.Errorf("TotalBlocks = %d, want 2", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 1 {
		t.Errorf("CoveredBlocks = %d, want 1", fc.CoveredBlocks)
	}
	if fc.CoveredPct != 50 {
		t.Errorf("CoveredPct = %f, want 50", fc.CoveredPct)
	}
	if len(fc.UncoveredLines) != 1 || fc.UncoveredLines[0] != 16 {
		t.Errorf("UncoveredLines = %v, want [16]", fc.UncoveredLines)
	}
}

func TestComputeFuncCoverageFullyCovered(t *testing.T) {
	blocks := []coverBlock{
		{file: "github.com/example/pkg/handler.go", startLine: 10, endLine: 15, count: 1},
		{file: "github.com/example/pkg/handler.go", startLine: 16, endLine: 20, count: 3},
	}

	r := funcRange{file: "pkg/handler.go", startLine: 10, endLine: 20, funcName: "Handler"}
	fc := computeFuncCoverage(r, blocks, "/project")

	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100", fc.CoveredPct)
	}
	if len(fc.UncoveredLines) != 0 {
		t.Errorf("UncoveredLines = %v, want []", fc.UncoveredLines)
	}
}

func TestComputeFuncCoverageNoMatchingBlocks(t *testing.T) {
	r := funcRange{file: "pkg/empty.go", startLine: 1, endLine: 5, funcName: "Empty"}
	fc := computeFuncCoverage(r, nil, "/project")

	if fc.CoveredPct != 0 {
		t.Errorf("CoveredPct = %f, want 0", fc.CoveredPct)
	}
	if fc.TotalBlocks != 0 {
		t.Errorf("TotalBlocks = %d, want 0", fc.TotalBlocks)
	}
}

func TestComputeFuncCoverageKey(t *testing.T) {
	r := funcRange{file: "pkg/handler.go", startLine: 10, endLine: 20, funcName: "Handler"}
	fc := computeFuncCoverage(r, nil, "/project")

	want := "pkg/handler.go:10-20"
	if fc.Key != want {
		t.Errorf("Key = %q, want %q", fc.Key, want)
	}
}
