package coverage

import "testing"

func TestCountJavaLines(t *testing.T) {
	fileCov := &jacocoFile{Lines: []jacocoLine{
		{Nr: 5, Ci: 4},        // covered
		{Nr: 6, Mi: 3, Ci: 0}, // uncovered (instructions but none covered)
		{Nr: 7},               // no instructions -> ignored
		{Nr: 99, Ci: 1},       // outside range -> ignored
	}}
	r := javaFuncRange{file: "Foo.java", startLine: 5, endLine: 10}
	var fc FuncCoverage
	countJavaLines(fileCov, r, &fc)

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
