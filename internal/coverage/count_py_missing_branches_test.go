package coverage

import "testing"

func TestCountPyMissingBranches(t *testing.T) {
	tests := []struct {
		name          string
		fileCov       *pyCoverageFile
		r             pyFuncRange
		wantTotal     int
		wantUncovered int
	}{
		{
			name: "missing branches in range",
			fileCov: &pyCoverageFile{
				MissingBranches: [][]int{{10, 15}, {12, 20}, {30, 35}},
			},
			r:             pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:     2,
			wantUncovered: 2,
		},
		{
			name: "no missing branches in range",
			fileCov: &pyCoverageFile{
				MissingBranches: [][]int{{30, 35}},
			},
			r:             pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:     0,
			wantUncovered: 0,
		},
		{
			name: "short branch entry skipped",
			fileCov: &pyCoverageFile{
				MissingBranches: [][]int{{10}, {12, 15}},
			},
			r:             pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:     1,
			wantUncovered: 1,
		},
		{
			name: "empty missing branches",
			fileCov: &pyCoverageFile{
				MissingBranches: nil,
			},
			r:             pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:     0,
			wantUncovered: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FuncCoverage{}
			countPyMissingBranches(tt.fileCov, tt.r, fc)
			if fc.TotalBlocks != tt.wantTotal {
				t.Errorf("TotalBlocks = %d, want %d", fc.TotalBlocks, tt.wantTotal)
			}
			if len(fc.UncoveredLines) != tt.wantUncovered {
				t.Errorf("UncoveredLines count = %d, want %d", len(fc.UncoveredLines), tt.wantUncovered)
			}
		})
	}
}
