package coverage

import "testing"

func TestComputeCsFuncCoverageNilCov(t *testing.T) {
	r := csFuncRange{file: "Foo.cs", startLine: 1, endLine: 5}
	fc := computeCsFuncCoverage(r, nil, "/proj")
	if fc.CoveredPct != 100 {
		t.Errorf("nil cov CoveredPct = %v, want 100", fc.CoveredPct)
	}
}

func TestComputeCsFuncCoverageFileNotFound(t *testing.T) {
	cov := &csCoverage{Files: []csFile{{Path: "App/Other.cs"}}}
	r := csFuncRange{file: "App/Foo.cs", startLine: 1, endLine: 5}
	fc := computeCsFuncCoverage(r, cov, "/proj")
	if fc.CoveredPct != 100 {
		t.Errorf("missing file CoveredPct = %v, want 100", fc.CoveredPct)
	}
}

func TestComputeCsFuncCoveragePartial(t *testing.T) {
	cov, err := parseCobertura(coberturaFixture)
	if err != nil {
		t.Fatal(err)
	}
	// Classify(): lines 9 (branch line, 1/2), 10 (covered), 12 (uncovered).
	r := csFuncRange{file: "App/Calculator.cs", startLine: 9, endLine: 12, funcName: "Classify"}
	fc := computeCsFuncCoverage(r, cov, "/proj")
	if fc.TotalBlocks == 0 {
		t.Fatal("expected some blocks counted")
	}
	if fc.CoveredPct >= 100 {
		t.Errorf("CoveredPct = %v, want < 100 for partial coverage", fc.CoveredPct)
	}
	if fc.Key != "App/Calculator.cs:9-12" {
		t.Errorf("Key = %q", fc.Key)
	}
}

func TestComputeCsFuncCoverageNoBlocks(t *testing.T) {
	cov, err := parseCobertura(coberturaFixture)
	if err != nil {
		t.Fatal(err)
	}
	r := csFuncRange{file: "App/Calculator.cs", startLine: 10000, endLine: 10001, funcName: "phantom"}
	fc := computeCsFuncCoverage(r, cov, "/proj")
	if fc.TotalBlocks != 0 {
		t.Fatalf("expected TotalBlocks=0 for empty range, got %d", fc.TotalBlocks)
	}
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %v, want 100 when no blocks", fc.CoveredPct)
	}
}

func TestComputeCsFuncCoverageFull(t *testing.T) {
	cov, err := parseCobertura(coberturaFixture)
	if err != nil {
		t.Fatal(err)
	}
	// Add(): line 5 only, fully covered, no branches.
	r := csFuncRange{file: "App/Calculator.cs", startLine: 5, endLine: 5, funcName: "Add"}
	fc := computeCsFuncCoverage(r, cov, "/proj")
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %v, want 100", fc.CoveredPct)
	}
}
