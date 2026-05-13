package coverage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// ---------------------------------------------------------------------------
// parsePyCoverageJSON tests
// ---------------------------------------------------------------------------

func TestParsePyCoverageJSON(t *testing.T) {
	dir := t.TempDir()

	covData := pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"handler.py": {
				ExecutedLines:    []int{1, 2, 3, 5, 6},
				MissingLines:    []int{4, 7},
				ExecutedBranches: [][]int{{3, 5}, {3, 7}},
				MissingBranches:  [][]int{{5, 6}},
			},
		},
	}

	raw, err := json.Marshal(covData)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "coverage.json")
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := parsePyCoverageJSON(path)
	if err != nil {
		t.Fatalf("parsePyCoverageJSON: %v", err)
	}

	fileCov, ok := result.Files["handler.py"]
	if !ok {
		t.Fatal("expected key 'handler.py' in result")
	}
	if len(fileCov.ExecutedLines) != 5 {
		t.Errorf("ExecutedLines length = %d, want 5", len(fileCov.ExecutedLines))
	}
	if len(fileCov.MissingLines) != 2 {
		t.Errorf("MissingLines length = %d, want 2", len(fileCov.MissingLines))
	}
	if len(fileCov.ExecutedBranches) != 2 {
		t.Errorf("ExecutedBranches length = %d, want 2", len(fileCov.ExecutedBranches))
	}
	if len(fileCov.MissingBranches) != 1 {
		t.Errorf("MissingBranches length = %d, want 1", len(fileCov.MissingBranches))
	}
}

func TestParsePyCoverageJSONFileNotFound(t *testing.T) {
	_, err := parsePyCoverageJSON("/nonexistent/coverage.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestParsePyCoverageJSONInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "coverage.json")
	if err := os.WriteFile(path, []byte("{invalid json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := parsePyCoverageJSON(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ---------------------------------------------------------------------------
// matchesPyPath tests
// ---------------------------------------------------------------------------

func TestMatchesPyPathDirect(t *testing.T) {
	if !matchesPyPath("handler.py", "handler.py", "/project") {
		t.Error("expected true for direct match")
	}
}

func TestMatchesPyPathAbsolute(t *testing.T) {
	if !matchesPyPath("/project/handler.py", "handler.py", "/project") {
		t.Error("expected true when covPath is absolute and matches projectRoot + targetFile")
	}
}

func TestMatchesPyPathSuffix(t *testing.T) {
	if !matchesPyPath("/some/deep/path/handler.py", "handler.py", "/project") {
		t.Error("expected true for suffix match")
	}
}

func TestMatchesPyPathNoMatch(t *testing.T) {
	if matchesPyPath("handler.py", "service.py", "/project") {
		t.Error("expected false for non-matching files")
	}
}

func TestMatchesPyPathSubdirectory(t *testing.T) {
	if !matchesPyPath("/project/src/handler.py", "src/handler.py", "/project") {
		t.Error("expected true for subdirectory match via absolute join")
	}
}

func TestMatchesPyPathDifferentDir(t *testing.T) {
	if matchesPyPath("other/handler.py", "src/handler.py", "/project") {
		t.Error("expected false when directories differ")
	}
}

// ---------------------------------------------------------------------------
// collectPyRanges tests
// ---------------------------------------------------------------------------

func TestCollectPyRanges(t *testing.T) {
	ep := &model.Endpoint{
		Name: "CreateOrder",
		Handler: model.FuncLocation{
			File:      "handlers/order.py",
			StartLine: 10,
			EndLine:   40,
		},
		Chain: []model.ChainEntry{
			{Func: "validate_order", File: "services/validator.py", StartLine: 5, EndLine: 25},
			{Func: "save_to_db", File: "repo/order.py", StartLine: 10, EndLine: 50, Boundary: "db"},
			{Func: "format_response", File: "helpers/format.py", StartLine: 1, EndLine: 15},
		},
	}

	ranges := collectPyRanges(ep)

	// Handler + 2 chain entries without boundary (save_to_db has boundary "db", excluded).
	if len(ranges) != 3 {
		t.Fatalf("got %d ranges, want 3", len(ranges))
	}

	if ranges[0].file != "handlers/order.py" || ranges[0].funcName != "CreateOrder" {
		t.Errorf("ranges[0] = {file: %q, func: %q}, want handler entry", ranges[0].file, ranges[0].funcName)
	}
	if ranges[0].startLine != 10 || ranges[0].endLine != 40 {
		t.Errorf("ranges[0] lines = %d-%d, want 10-40", ranges[0].startLine, ranges[0].endLine)
	}

	if ranges[1].file != "services/validator.py" || ranges[1].funcName != "validate_order" {
		t.Errorf("ranges[1] = {file: %q, func: %q}, want validate_order", ranges[1].file, ranges[1].funcName)
	}

	if ranges[2].file != "helpers/format.py" || ranges[2].funcName != "format_response" {
		t.Errorf("ranges[2] = {file: %q, func: %q}, want format_response", ranges[2].file, ranges[2].funcName)
	}
}

func TestCollectPyRangesEmptyEndpoint(t *testing.T) {
	ep := &model.Endpoint{
		Name:    "Empty",
		Handler: model.FuncLocation{},
		Chain:   nil,
	}

	ranges := collectPyRanges(ep)
	if len(ranges) != 0 {
		t.Errorf("expected empty ranges, got %d", len(ranges))
	}
}

// ---------------------------------------------------------------------------
// computePyFuncCoverage tests
// ---------------------------------------------------------------------------

func TestComputePyFuncCoverage(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"handler.py": {
				ExecutedLines:    []int{10, 11, 12, 15},
				MissingLines:    []int{13, 14},
				ExecutedBranches: [][]int{{12, 15}},
				MissingBranches:  [][]int{{12, 13}},
			},
		},
	}

	r := pyFuncRange{file: "handler.py", startLine: 10, endLine: 20, funcName: "handler"}
	fc := computePyFuncCoverage(r, covData, "/project")

	// Executed lines in range: 10, 11, 12, 15 = 4 covered
	// Missing lines in range: 13, 14 = 2 uncovered
	// Executed branches in range: [12,15] where 12 is in range = 1 covered
	// Missing branches in range: [12,13] where 12 is in range = 1 uncovered
	// Total = 4 + 2 + 1 + 1 = 8
	// Covered = 4 + 1 = 5
	if fc.TotalBlocks != 8 {
		t.Errorf("TotalBlocks = %d, want 8", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 5 {
		t.Errorf("CoveredBlocks = %d, want 5", fc.CoveredBlocks)
	}
	if fc.CoveredPct != 62.5 {
		t.Errorf("CoveredPct = %f, want 62.5", fc.CoveredPct)
	}
	if len(fc.UncoveredLines) != 3 {
		t.Errorf("UncoveredLines = %v, want 3 entries", fc.UncoveredLines)
	}
}

func TestComputePyFuncCoverageNilData(t *testing.T) {
	r := pyFuncRange{file: "handler.py", startLine: 1, endLine: 10, funcName: "handler"}
	fc := computePyFuncCoverage(r, nil, "/project")
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100 (nil data)", fc.CoveredPct)
	}
}

func TestComputePyFuncCoverageFileNotInData(t *testing.T) {
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"other.py": {ExecutedLines: []int{1, 2}, MissingLines: nil},
		},
	}

	r := pyFuncRange{file: "handler.py", startLine: 1, endLine: 10, funcName: "handler"}
	fc := computePyFuncCoverage(r, covData, "/project")
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100 (file not in coverage data)", fc.CoveredPct)
	}
}

func TestComputePyFuncCoverageOutOfRange(t *testing.T) {
	// All lines are outside the function range.
	covData := &pyCoverageJSON{
		Files: map[string]pyCoverageFile{
			"handler.py": {
				ExecutedLines: []int{50, 60},
				MissingLines:  []int{70, 80},
			},
		},
	}

	r := pyFuncRange{file: "handler.py", startLine: 1, endLine: 10, funcName: "handler"}
	fc := computePyFuncCoverage(r, covData, "/project")
	// No blocks in range, treated as 100%.
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100 (no lines in range)", fc.CoveredPct)
	}
}
