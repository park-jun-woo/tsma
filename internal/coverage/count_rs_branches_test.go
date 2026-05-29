package coverage

import "testing"

func TestCountRsBranches(t *testing.T) {
	cov, err := parseLLVMCov(llvmCovFixture)
	if err != nil {
		t.Fatal(err)
	}
	fileCov := findRsCoverageFile(cov, "src/lib.rs", "/home/dev/demo")
	if fileCov == nil {
		t.Fatal("no file cov")
	}

	// classify branch at line 3: both sides covered (exec 2, false 3).
	var fc FuncCoverage
	countRsBranches(fileCov, rsFuncRange{file: "src/lib.rs", startLine: 1, endLine: 5}, &fc)
	if fc.TotalBlocks != 2 || fc.CoveredBlocks != 2 {
		t.Errorf("classify branches: total=%d covered=%d, want 2/2", fc.TotalBlocks, fc.CoveredBlocks)
	}

	// unused branch at line 9: true side uncovered, false side covered.
	var fc2 FuncCoverage
	countRsBranches(fileCov, rsFuncRange{file: "src/lib.rs", startLine: 8, endLine: 10}, &fc2)
	if fc2.TotalBlocks != 2 || fc2.CoveredBlocks != 1 {
		t.Errorf("unused branches: total=%d covered=%d, want 2/1", fc2.TotalBlocks, fc2.CoveredBlocks)
	}
}

func TestCountRsBranchesFalseSideUncovered(t *testing.T) {
	// Synthetic branch: true side covered (exec>0), false side uncovered
	// (FalseExecCount==0) so both the covered-true and uncovered-false paths run.
	fileCov := &llvmCovFile{
		Filename: "/home/dev/demo/src/lib.rs",
		Branches: []llvmBranch{
			{LineStart: 12, ExecCount: 5, FalseExecCount: 0},
			// Out-of-range branch must be skipped.
			{LineStart: 999, ExecCount: 0, FalseExecCount: 0},
		},
	}

	var fc FuncCoverage
	countRsBranches(fileCov, rsFuncRange{file: "src/lib.rs", startLine: 10, endLine: 15}, &fc)

	if fc.TotalBlocks != 2 {
		t.Errorf("TotalBlocks = %d, want 2 (one in-range branch, two sides)", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 1 {
		t.Errorf("CoveredBlocks = %d, want 1 (only the true side covered)", fc.CoveredBlocks)
	}
	// The uncovered false side records the branch line.
	if len(fc.UncoveredLines) != 1 || fc.UncoveredLines[0] != 12 {
		t.Errorf("UncoveredLines = %v, want [12]", fc.UncoveredLines)
	}
}

func TestCountRsBranchesBothSidesUncovered(t *testing.T) {
	// Both sides uncovered -> both UncoveredLines appends run.
	fileCov := &llvmCovFile{
		Filename: "/home/dev/demo/src/lib.rs",
		Branches: []llvmBranch{
			{LineStart: 20, ExecCount: 0, FalseExecCount: 0},
		},
	}

	var fc FuncCoverage
	countRsBranches(fileCov, rsFuncRange{file: "src/lib.rs", startLine: 18, endLine: 22}, &fc)

	if fc.TotalBlocks != 2 || fc.CoveredBlocks != 0 {
		t.Errorf("total=%d covered=%d, want 2/0", fc.TotalBlocks, fc.CoveredBlocks)
	}
	if len(fc.UncoveredLines) != 2 {
		t.Errorf("UncoveredLines = %v, want both sides recorded", fc.UncoveredLines)
	}
}
