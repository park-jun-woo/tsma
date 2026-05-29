package coverage

import "testing"

func TestCountCsBranches(t *testing.T) {
	fileCov := &csFile{Lines: []coberturaLine{
		{Number: 9, Branch: "true", ConditionCoverage: "50% (1/2)"},   // 1 covered, 1 missed
		{Number: 10, Branch: "true", ConditionCoverage: "100% (2/2)"}, // both covered
		{Number: 11, Branch: "false"},                                 // not a branch -> ignored
		{Number: 12, Branch: "true", ConditionCoverage: ""},           // branch but malformed/empty cc -> total==0 skipped
		{Number: 50, Branch: "true", ConditionCoverage: "0% (0/4)"},   // outside range -> ignored
	}}
	r := csFuncRange{file: "Foo.cs", startLine: 9, endLine: 12}
	var fc FuncCoverage
	countCsBranches(fileCov, r, &fc)

	if fc.TotalBlocks != 4 {
		t.Errorf("TotalBlocks = %d, want 4", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 3 {
		t.Errorf("CoveredBlocks = %d, want 3", fc.CoveredBlocks)
	}
	if len(fc.UncoveredLines) != 1 || fc.UncoveredLines[0] != 9 {
		t.Errorf("UncoveredLines = %v, want [9]", fc.UncoveredLines)
	}
}
