package coverage

import "testing"

func TestComputeTSFuncCoverageWithData(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"src/handler.ts": {
			StatementMap: map[string]coverageRange{
				"0": {Start: coveragePosition{Line: 10}, End: coveragePosition{Line: 10}},
				"1": {Start: coveragePosition{Line: 15}, End: coveragePosition{Line: 15}},
			},
			S:         map[string]int{"0": 1, "1": 0},
			BranchMap: map[string]coverageBranch{},
			B:         map[string][]int{},
			FnMap:     map[string]coverageFunction{},
			F:         map[string]int{},
		},
	}

	r := tsFuncRange{file: "src/handler.ts", startLine: 10, endLine: 20, funcName: "handler"}
	fc := computeTSFuncCoverage(r, data, "/project")

	if fc.TotalBlocks != 2 {
		t.Errorf("TotalBlocks = %d, want 2", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 1 {
		t.Errorf("CoveredBlocks = %d, want 1", fc.CoveredBlocks)
	}
	if fc.CoveredPct != 50 {
		t.Errorf("CoveredPct = %f, want 50", fc.CoveredPct)
	}
}

func TestComputeTSFuncCoverageNoEntry(t *testing.T) {
	data := map[string]coverageFinalEntry{}
	r := tsFuncRange{file: "src/missing.ts", startLine: 1, endLine: 10, funcName: "missing"}
	fc := computeTSFuncCoverage(r, data, "/project")
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100 (no entry)", fc.CoveredPct)
	}
}

func TestComputeTSFuncCoverageEmptyStatements(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"src/handler.ts": {
			StatementMap: map[string]coverageRange{
				"0": {Start: coveragePosition{Line: 50}, End: coveragePosition{Line: 50}},
			},
			S:         map[string]int{"0": 1},
			BranchMap: map[string]coverageBranch{},
			B:         map[string][]int{},
			FnMap:     map[string]coverageFunction{},
			F:         map[string]int{},
		},
	}

	r := tsFuncRange{file: "src/handler.ts", startLine: 1, endLine: 10, funcName: "empty"}
	fc := computeTSFuncCoverage(r, data, "/project")
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100 (no stmts in range)", fc.CoveredPct)
	}
}

func TestComputeTSFuncCoverageKey(t *testing.T) {
	data := map[string]coverageFinalEntry{}
	r := tsFuncRange{file: "src/handler.ts", startLine: 5, endLine: 15, funcName: "handler"}
	fc := computeTSFuncCoverage(r, data, "/project")
	want := "src/handler.ts:5-15"
	if fc.Key != want {
		t.Errorf("Key = %q, want %q", fc.Key, want)
	}
}
