package coverage

import "testing"

func TestCountRsLineSegments(t *testing.T) {
	cov, err := parseLLVMCov(llvmCovFixture)
	if err != nil {
		t.Fatal(err)
	}
	fileCov := findRsCoverageFile(cov, "src/lib.rs", "/home/dev/demo")
	if fileCov == nil {
		t.Fatal("no file cov")
	}

	// classify: lines 1-5, all executed region entries (4 lines), all covered.
	var fc FuncCoverage
	countRsLineSegments(fileCov, rsFuncRange{file: "src/lib.rs", startLine: 1, endLine: 5}, &fc)
	if fc.TotalBlocks != 4 || fc.CoveredBlocks != 4 {
		t.Errorf("classify lines: total=%d covered=%d, want 4/4", fc.TotalBlocks, fc.CoveredBlocks)
	}
	if len(fc.UncoveredLines) != 0 {
		t.Errorf("classify uncovered lines = %v, want none", fc.UncoveredLines)
	}

	// unused: lines 8-10, line 8 covered, line 9 not.
	var fc2 FuncCoverage
	countRsLineSegments(fileCov, rsFuncRange{file: "src/lib.rs", startLine: 8, endLine: 10}, &fc2)
	if fc2.TotalBlocks != 2 || fc2.CoveredBlocks != 1 {
		t.Errorf("unused lines: total=%d covered=%d, want 2/1", fc2.TotalBlocks, fc2.CoveredBlocks)
	}
	if len(fc2.UncoveredLines) != 1 || fc2.UncoveredLines[0] != 9 {
		t.Errorf("unused uncovered lines = %v, want [9]", fc2.UncoveredLines)
	}
}
