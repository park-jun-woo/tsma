package coverage

import "testing"

func TestCountPyExecutedBranches(t *testing.T) {
	tests := []struct {
		name           string
		fileCov        *pyCoverageFile
		r              pyFuncRange
		wantTotal      int
		wantCovered    int
	}{
		{
			name: "branches in range",
			fileCov: &pyCoverageFile{
				ExecutedBranches: [][]int{{10, 15}, {12, 20}, {30, 35}},
			},
			r:           pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:   2,
			wantCovered: 2,
		},
		{
			name: "no branches in range",
			fileCov: &pyCoverageFile{
				ExecutedBranches: [][]int{{30, 35}, {40, 45}},
			},
			r:           pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:   0,
			wantCovered: 0,
		},
		{
			name: "empty branches",
			fileCov: &pyCoverageFile{
				ExecutedBranches: nil,
			},
			r:           pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:   0,
			wantCovered: 0,
		},
		{
			name: "short branch entry skipped",
			fileCov: &pyCoverageFile{
				ExecutedBranches: [][]int{{10}, {12, 15}},
			},
			r:           pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:   1,
			wantCovered: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FuncCoverage{}
			countPyExecutedBranches(tt.fileCov, tt.r, fc)
			if fc.TotalBlocks != tt.wantTotal {
				t.Errorf("TotalBlocks = %d, want %d", fc.TotalBlocks, tt.wantTotal)
			}
			if fc.CoveredBlocks != tt.wantCovered {
				t.Errorf("CoveredBlocks = %d, want %d", fc.CoveredBlocks, tt.wantCovered)
			}
		})
	}
}
