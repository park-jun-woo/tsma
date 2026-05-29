package coverage

import "testing"

func TestComputeRsFuncCoverageFullyCovered(t *testing.T) {
	cov, err := parseLLVMCov(llvmCovFixture)
	if err != nil {
		t.Fatal(err)
	}
	r := rsFuncRange{file: "src/lib.rs", startLine: 1, endLine: 5, funcName: "classify"}
	fc := computeRsFuncCoverage(r, cov, "/home/dev/demo")

	if fc.TotalBlocks != 6 || fc.CoveredBlocks != 6 {
		t.Errorf("classify: total=%d covered=%d, want 6/6", fc.TotalBlocks, fc.CoveredBlocks)
	}
	if fc.CoveredPct != 100 {
		t.Errorf("classify pct = %v, want 100", fc.CoveredPct)
	}
}

func TestComputeRsFuncCoveragePartial(t *testing.T) {
	cov, err := parseLLVMCov(llvmCovFixture)
	if err != nil {
		t.Fatal(err)
	}
	r := rsFuncRange{file: "src/lib.rs", startLine: 8, endLine: 10, funcName: "unused"}
	fc := computeRsFuncCoverage(r, cov, "/home/dev/demo")

	if fc.TotalBlocks != 4 || fc.CoveredBlocks != 2 {
		t.Errorf("unused: total=%d covered=%d, want 4/2", fc.TotalBlocks, fc.CoveredBlocks)
	}
	if fc.CoveredPct != 50 {
		t.Errorf("unused pct = %v, want 50", fc.CoveredPct)
	}
}

func TestComputeRsFuncCoverageNilData(t *testing.T) {
	fc := computeRsFuncCoverage(rsFuncRange{file: "src/lib.rs", startLine: 1, endLine: 5}, nil, "/x")
	if fc.CoveredPct != 100 {
		t.Errorf("nil cov pct = %v, want 100", fc.CoveredPct)
	}
}

func TestComputeRsFuncCoverageFileFoundButNoSegments(t *testing.T) {
	cov, err := parseLLVMCov(llvmCovFixture)
	if err != nil {
		t.Fatal(err)
	}
	// File is present in the coverage data, but this line range contains no
	// executable segments or branches -> TotalBlocks stays 0 -> the else branch
	// sets CoveredPct to 100.
	r := rsFuncRange{file: "src/lib.rs", startLine: 10000, endLine: 10005, funcName: "phantom"}
	fc := computeRsFuncCoverage(r, cov, "/home/dev/demo")

	if fc.TotalBlocks != 0 {
		t.Errorf("expected TotalBlocks=0 for empty range, got %d", fc.TotalBlocks)
	}
	if fc.CoveredPct != 100 {
		t.Errorf("expected CoveredPct=100 for zero-block range, got %v", fc.CoveredPct)
	}
}

func TestComputeRsFuncCoverageUnknownFile(t *testing.T) {
	cov, err := parseLLVMCov(llvmCovFixture)
	if err != nil {
		t.Fatal(err)
	}
	fc := computeRsFuncCoverage(rsFuncRange{file: "src/missing.rs", startLine: 1, endLine: 5}, cov, "/home/dev/demo")
	if fc.CoveredPct != 100 {
		t.Errorf("unknown file pct = %v, want 100", fc.CoveredPct)
	}
}
