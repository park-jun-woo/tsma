package coverage

import "testing"

func TestComputeJavaFuncCoverageNilCov(t *testing.T) {
	r := javaFuncRange{file: "Foo.java", startLine: 1, endLine: 5}
	fc := computeJavaFuncCoverage(r, nil, "/proj")
	if fc.CoveredPct != 100 {
		t.Errorf("nil cov CoveredPct = %v, want 100", fc.CoveredPct)
	}
}

func TestComputeJavaFuncCoverageFileNotFound(t *testing.T) {
	cov := &jacocoCoverage{Files: []jacocoFile{{Path: "com/other/X.java"}}}
	r := javaFuncRange{file: "src/main/java/com/example/Foo.java", startLine: 1, endLine: 5}
	fc := computeJavaFuncCoverage(r, cov, "/proj")
	if fc.CoveredPct != 100 {
		t.Errorf("missing file CoveredPct = %v, want 100", fc.CoveredPct)
	}
}

func TestComputeJavaFuncCoveragePartial(t *testing.T) {
	cov, err := parseJacoco(jacocoFixture)
	if err != nil {
		t.Fatal(err)
	}
	// classify(): lines 9 (branch line), 10, 12; one branch missed and line 12 missed.
	r := javaFuncRange{file: "src/main/java/com/example/Calculator.java", startLine: 9, endLine: 12, funcName: "classify"}
	fc := computeJavaFuncCoverage(r, cov, "/proj")
	if fc.TotalBlocks == 0 {
		t.Fatal("expected some blocks counted")
	}
	if fc.CoveredPct >= 100 {
		t.Errorf("CoveredPct = %v, want < 100 for partial coverage", fc.CoveredPct)
	}
	if fc.Key != "src/main/java/com/example/Calculator.java:9-12" {
		t.Errorf("Key = %q", fc.Key)
	}
}

// TestComputeJavaFuncCoverageNoBlocks covers the TotalBlocks == 0 else-branch
// (line 33-34): the coverage file is found, but the function range covers no
// recorded executable or branch lines, so coverage defaults to 100.
func TestComputeJavaFuncCoverageNoBlocks(t *testing.T) {
	cov, err := parseJacoco(jacocoFixture)
	if err != nil {
		t.Fatal(err)
	}
	// A line range far beyond any recorded line in the fixture: the file is
	// matched but no line/branch counters fall inside the range.
	r := javaFuncRange{
		file:      "src/main/java/com/example/Calculator.java",
		startLine: 10000,
		endLine:   10001,
		funcName:  "phantom",
	}
	fc := computeJavaFuncCoverage(r, cov, "/proj")
	if fc.TotalBlocks != 0 {
		t.Fatalf("expected TotalBlocks=0 for empty range, got %d", fc.TotalBlocks)
	}
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %v, want 100 when no blocks", fc.CoveredPct)
	}
}

func TestComputeJavaFuncCoverageFull(t *testing.T) {
	cov, err := parseJacoco(jacocoFixture)
	if err != nil {
		t.Fatal(err)
	}
	// add(): line 5 only, fully covered, no branches.
	r := javaFuncRange{file: "src/main/java/com/example/Calculator.java", startLine: 5, endLine: 5, funcName: "add"}
	fc := computeJavaFuncCoverage(r, cov, "/proj")
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %v, want 100", fc.CoveredPct)
	}
}
