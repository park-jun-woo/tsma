package coverage

import "testing"

func TestCountTSBranchLocations(t *testing.T) {
	tests := []struct {
		name        string
		entry       *coverageFinalEntry
		branchID    string
		branch      coverageBranch
		r           tsFuncRange
		wantTotal   int
		wantCovered int
		wantUncov   int
	}{
		{
			name: "locations in range with mixed coverage",
			entry: &coverageFinalEntry{
				B: map[string][]int{"0": {1, 0, 1}},
			},
			branchID: "0",
			branch: coverageBranch{
				Locations: []coverageRange{
					{Start: coveragePosition{Line: 10}},
					{Start: coveragePosition{Line: 12}},
					{Start: coveragePosition{Line: 15}},
				},
			},
			r:           tsFuncRange{startLine: 10, endLine: 20},
			wantTotal:   3,
			wantCovered: 2,
			wantUncov:   1,
		},
		{
			name: "locations out of range",
			entry: &coverageFinalEntry{
				B: map[string][]int{"0": {1}},
			},
			branchID: "0",
			branch: coverageBranch{
				Locations: []coverageRange{
					{Start: coveragePosition{Line: 30}},
				},
			},
			r:           tsFuncRange{startLine: 10, endLine: 20},
			wantTotal:   0,
			wantCovered: 0,
			wantUncov:   0,
		},
		{
			name: "count index out of range treated as uncovered",
			entry: &coverageFinalEntry{
				B: map[string][]int{"0": {1}}, // only 1 count but 2 locations
			},
			branchID: "0",
			branch: coverageBranch{
				Locations: []coverageRange{
					{Start: coveragePosition{Line: 10}},
					{Start: coveragePosition{Line: 12}},
				},
			},
			r:           tsFuncRange{startLine: 10, endLine: 20},
			wantTotal:   2,
			wantCovered: 1,
			wantUncov:   1,
		},
		{
			name: "empty locations",
			entry: &coverageFinalEntry{
				B: map[string][]int{"0": {}},
			},
			branchID: "0",
			branch: coverageBranch{
				Locations: nil,
			},
			r:           tsFuncRange{startLine: 10, endLine: 20},
			wantTotal:   0,
			wantCovered: 0,
			wantUncov:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FuncCoverage{}
			countTSBranchLocations(tt.entry, tt.branchID, tt.branch, tt.r, fc)
			if fc.TotalBlocks != tt.wantTotal {
				t.Errorf("TotalBlocks = %d, want %d", fc.TotalBlocks, tt.wantTotal)
			}
			if fc.CoveredBlocks != tt.wantCovered {
				t.Errorf("CoveredBlocks = %d, want %d", fc.CoveredBlocks, tt.wantCovered)
			}
			if len(fc.UncoveredLines) != tt.wantUncov {
				t.Errorf("UncoveredLines = %d, want %d", len(fc.UncoveredLines), tt.wantUncov)
			}
		})
	}
}
