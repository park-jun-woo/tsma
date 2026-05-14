package coverage

import "testing"

func TestCountTSBranches(t *testing.T) {
	entry := &coverageFinalEntry{
		BranchMap: map[string]coverageBranch{
			"0": {
				Loc: coverageRange{Start: coveragePosition{Line: 10}},
				Locations: []coverageRange{
					{Start: coveragePosition{Line: 10}},
					{Start: coveragePosition{Line: 11}},
				},
			},
			"1": {
				Loc: coverageRange{Start: coveragePosition{Line: 50}}, // outside range
				Locations: []coverageRange{
					{Start: coveragePosition{Line: 50}},
				},
			},
		},
		B: map[string][]int{
			"0": {1, 0},
			"1": {1},
		},
	}

	r := tsFuncRange{startLine: 5, endLine: 20}
	fc := &FuncCoverage{}
	countTSBranches(entry, r, fc)

	if fc.TotalBlocks != 2 {
		t.Errorf("TotalBlocks = %d, want 2", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 1 {
		t.Errorf("CoveredBlocks = %d, want 1", fc.CoveredBlocks)
	}
}

func TestCountTSBranchesEmpty(t *testing.T) {
	entry := &coverageFinalEntry{
		BranchMap: map[string]coverageBranch{},
		B:         map[string][]int{},
	}

	r := tsFuncRange{startLine: 1, endLine: 10}
	fc := &FuncCoverage{}
	countTSBranches(entry, r, fc)

	if fc.TotalBlocks != 0 {
		t.Errorf("TotalBlocks = %d, want 0", fc.TotalBlocks)
	}
}
