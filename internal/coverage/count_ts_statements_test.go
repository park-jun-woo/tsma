package coverage

import "testing"

func TestCountTSStatements(t *testing.T) {
	entry := &coverageFinalEntry{
		StatementMap: map[string]coverageRange{
			"0": {Start: coveragePosition{Line: 10}, End: coveragePosition{Line: 10}},
			"1": {Start: coveragePosition{Line: 15}, End: coveragePosition{Line: 15}},
			"2": {Start: coveragePosition{Line: 50}, End: coveragePosition{Line: 50}}, // outside
		},
		S: map[string]int{"0": 1, "1": 0, "2": 1},
	}

	r := tsFuncRange{startLine: 5, endLine: 20}
	fc := &FuncCoverage{}
	countTSStatements(entry, r, fc)

	if fc.TotalBlocks != 2 {
		t.Errorf("TotalBlocks = %d, want 2", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 1 {
		t.Errorf("CoveredBlocks = %d, want 1", fc.CoveredBlocks)
	}
	if len(fc.UncoveredLines) != 1 || fc.UncoveredLines[0] != 15 {
		t.Errorf("UncoveredLines = %v, want [15]", fc.UncoveredLines)
	}
}

func TestCountTSStatementsEmpty(t *testing.T) {
	entry := &coverageFinalEntry{
		StatementMap: map[string]coverageRange{},
		S:            map[string]int{},
	}

	r := tsFuncRange{startLine: 1, endLine: 10}
	fc := &FuncCoverage{}
	countTSStatements(entry, r, fc)

	if fc.TotalBlocks != 0 {
		t.Errorf("TotalBlocks = %d, want 0", fc.TotalBlocks)
	}
}

func TestCountTSStatementsAllCovered(t *testing.T) {
	entry := &coverageFinalEntry{
		StatementMap: map[string]coverageRange{
			"0": {Start: coveragePosition{Line: 10}, End: coveragePosition{Line: 10}},
			"1": {Start: coveragePosition{Line: 12}, End: coveragePosition{Line: 12}},
		},
		S: map[string]int{"0": 3, "1": 1},
	}

	r := tsFuncRange{startLine: 10, endLine: 20}
	fc := &FuncCoverage{}
	countTSStatements(entry, r, fc)

	if fc.TotalBlocks != 2 {
		t.Errorf("TotalBlocks = %d, want 2", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 2 {
		t.Errorf("CoveredBlocks = %d, want 2", fc.CoveredBlocks)
	}
	if len(fc.UncoveredLines) != 0 {
		t.Errorf("UncoveredLines = %v, want []", fc.UncoveredLines)
	}
}
