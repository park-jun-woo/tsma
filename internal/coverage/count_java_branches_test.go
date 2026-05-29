package coverage

import "testing"

func TestCountJavaBranches(t *testing.T) {
	fileCov := &jacocoFile{Lines: []jacocoLine{
		{Nr: 9, Cb: 1, Mb: 1},  // one covered, one missed
		{Nr: 10, Cb: 2, Mb: 0}, // both covered
		{Nr: 50, Cb: 0, Mb: 4}, // outside range -> ignored
	}}
	r := javaFuncRange{file: "Foo.java", startLine: 9, endLine: 12}
	var fc FuncCoverage
	countJavaBranches(fileCov, r, &fc)

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
