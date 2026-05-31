package coverage

import "testing"

func TestComputePyFuncCoverageWithData(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"handler.py": {
				ExecutedLines:    []int{10, 11},
				MissingLines:     []int{12},
				ExecutedBranches: [][]int{{11, 15}},
				MissingBranches:  nil,
			},
		},
	}

	r := pyFuncRange{file: "handler.py", startLine: 10, endLine: 20, funcName: "handler"}
	fc := computePyFuncCoverage(r, covData, "/project")

	// 2 executed lines + 1 missing line + 1 executed branch = 4 total, 3 covered
	if fc.TotalBlocks != 4 {
		t.Errorf("TotalBlocks = %d, want 4", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 3 {
		t.Errorf("CoveredBlocks = %d, want 3", fc.CoveredBlocks)
	}
	if fc.CoveredPct != 75 {
		t.Errorf("CoveredPct = %f, want 75", fc.CoveredPct)
	}
}

func TestComputePyFuncCoverageNilCovData(t *testing.T) {
	r := pyFuncRange{file: "handler.py", startLine: 1, endLine: 10, funcName: "handler"}
	fc := computePyFuncCoverage(r, nil, "/project")
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100 (nil data)", fc.CoveredPct)
	}
}

func TestComputePyFuncCoverageFileNotFound(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"other.py": {ExecutedLines: []int{1}},
		},
	}
	r := pyFuncRange{file: "handler.py", startLine: 1, endLine: 10, funcName: "handler"}
	fc := computePyFuncCoverage(r, covData, "/project")
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100 (file not found)", fc.CoveredPct)
	}
}

func TestComputePyFuncCoverageNoBlocksInRange(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"handler.py": {
				ExecutedLines: []int{50, 60},
				MissingLines:  []int{70},
			},
		},
	}
	r := pyFuncRange{file: "handler.py", startLine: 1, endLine: 10, funcName: "handler"}
	fc := computePyFuncCoverage(r, covData, "/project")
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100 (no blocks in range)", fc.CoveredPct)
	}
}

func TestComputePyFuncCoverageKey(t *testing.T) {
	r := pyFuncRange{file: "handler.py", startLine: 5, endLine: 15, funcName: "handler"}
	fc := computePyFuncCoverage(r, nil, "/project")
	want := "handler.py:5-15"
	if fc.Key != want {
		t.Errorf("Key = %q, want %q", fc.Key, want)
	}
}
