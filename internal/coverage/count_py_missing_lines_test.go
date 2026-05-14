package coverage

import "testing"

func TestCountPyMissingLines(t *testing.T) {
	tests := []struct {
		name          string
		fileCov       *pyCoverageFile
		r             pyFuncRange
		wantTotal     int
		wantUncovered int
	}{
		{
			name: "missing lines in range",
			fileCov: &pyCoverageFile{
				MissingLines: []int{10, 15, 25},
			},
			r:             pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:     2,
			wantUncovered: 2,
		},
		{
			name: "no missing lines in range",
			fileCov: &pyCoverageFile{
				MissingLines: []int{1, 2, 30},
			},
			r:             pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:     0,
			wantUncovered: 0,
		},
		{
			name: "empty missing lines",
			fileCov: &pyCoverageFile{
				MissingLines: nil,
			},
			r:             pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:     0,
			wantUncovered: 0,
		},
		{
			name: "boundary lines included",
			fileCov: &pyCoverageFile{
				MissingLines: []int{10, 20},
			},
			r:             pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:     2,
			wantUncovered: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FuncCoverage{}
			countPyMissingLines(tt.fileCov, tt.r, fc)
			if fc.TotalBlocks != tt.wantTotal {
				t.Errorf("TotalBlocks = %d, want %d", fc.TotalBlocks, tt.wantTotal)
			}
			if len(fc.UncoveredLines) != tt.wantUncovered {
				t.Errorf("UncoveredLines count = %d, want %d", len(fc.UncoveredLines), tt.wantUncovered)
			}
		})
	}
}
