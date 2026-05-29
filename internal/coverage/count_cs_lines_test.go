package coverage

import "testing"

func TestCountCsLines(t *testing.T) {
	fileCov := &csFile{Lines: []coberturaLine{
		{Number: 5, Hits: 1},  // covered
		{Number: 6, Hits: 0},  // uncovered
		{Number: 99, Hits: 1}, // outside range -> ignored
	}}
	r := csFuncRange{file: "Foo.cs", startLine: 5, endLine: 10}
	var fc FuncCoverage
	countCsLines(fileCov, r, &fc)

	if fc.TotalBlocks != 2 {
		t.Errorf("TotalBlocks = %d, want 2", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 1 {
		t.Errorf("CoveredBlocks = %d, want 1", fc.CoveredBlocks)
	}
	if len(fc.UncoveredLines) != 1 || fc.UncoveredLines[0] != 6 {
		t.Errorf("UncoveredLines = %v, want [6]", fc.UncoveredLines)
	}
}
