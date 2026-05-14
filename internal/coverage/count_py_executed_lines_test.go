package coverage

import "testing"

func TestCountPyExecutedLines(t *testing.T) {
	tests := []struct {
		name        string
		fileCov     *pyCoverageFile
		r           pyFuncRange
		wantTotal   int
		wantCovered int
	}{
		{
			name: "lines in range",
			fileCov: &pyCoverageFile{
				ExecutedLines: []int{10, 11, 12, 25, 30},
			},
			r:           pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:   3,
			wantCovered: 3,
		},
		{
			name: "no lines in range",
			fileCov: &pyCoverageFile{
				ExecutedLines: []int{1, 2, 30},
			},
			r:           pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:   0,
			wantCovered: 0,
		},
		{
			name: "empty executed lines",
			fileCov: &pyCoverageFile{
				ExecutedLines: nil,
			},
			r:           pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:   0,
			wantCovered: 0,
		},
		{
			name: "boundary lines included",
			fileCov: &pyCoverageFile{
				ExecutedLines: []int{10, 20},
			},
			r:           pyFuncRange{startLine: 10, endLine: 20},
			wantTotal:   2,
			wantCovered: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FuncCoverage{}
			countPyExecutedLines(tt.fileCov, tt.r, fc)
			if fc.TotalBlocks != tt.wantTotal {
				t.Errorf("TotalBlocks = %d, want %d", fc.TotalBlocks, tt.wantTotal)
			}
			if fc.CoveredBlocks != tt.wantCovered {
				t.Errorf("CoveredBlocks = %d, want %d", fc.CoveredBlocks, tt.wantCovered)
			}
		})
	}
}
