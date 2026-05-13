package coverage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// ---------------------------------------------------------------------------
// parseCoverageFinalJSON tests
// ---------------------------------------------------------------------------

func TestParseCoverageFinalJSON(t *testing.T) {
	dir := t.TempDir()
	coverDir := filepath.Join(dir, ".tsma", "coverage")
	if err := os.MkdirAll(coverDir, 0o755); err != nil {
		t.Fatal(err)
	}

	data := map[string]coverageFinalEntry{
		"src/handler.ts": {
			StatementMap: map[string]coverageRange{
				"0": {
					Start: coveragePosition{Line: 1, Column: 0},
					End:   coveragePosition{Line: 1, Column: 20},
				},
				"1": {
					Start: coveragePosition{Line: 3, Column: 0},
					End:   coveragePosition{Line: 3, Column: 30},
				},
			},
			S: map[string]int{"0": 1, "1": 0},
			BranchMap: map[string]coverageBranch{
				"0": {
					Loc:  coverageRange{Start: coveragePosition{Line: 5, Column: 0}, End: coveragePosition{Line: 5, Column: 30}},
					Type: "if",
					Locations: []coverageRange{
						{Start: coveragePosition{Line: 5, Column: 0}, End: coveragePosition{Line: 5, Column: 15}},
						{Start: coveragePosition{Line: 5, Column: 16}, End: coveragePosition{Line: 5, Column: 30}},
					},
				},
			},
			B:     map[string][]int{"0": {1, 0}},
			FnMap: map[string]coverageFunction{},
			F:     map[string]int{},
		},
	}

	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(coverDir, "coverage-final.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := parseCoverageFinalJSON(coverDir)
	if err != nil {
		t.Fatalf("parseCoverageFinalJSON: %v", err)
	}

	entry, ok := result["src/handler.ts"]
	if !ok {
		t.Fatal("expected key 'src/handler.ts' in result")
	}

	if len(entry.StatementMap) != 2 {
		t.Errorf("StatementMap length = %d, want 2", len(entry.StatementMap))
	}
	if entry.S["0"] != 1 {
		t.Errorf("S[\"0\"] = %d, want 1", entry.S["0"])
	}
	if entry.S["1"] != 0 {
		t.Errorf("S[\"1\"] = %d, want 0", entry.S["1"])
	}
	if len(entry.BranchMap) != 1 {
		t.Errorf("BranchMap length = %d, want 1", len(entry.BranchMap))
	}
	branch := entry.BranchMap["0"]
	if branch.Type != "if" {
		t.Errorf("branch type = %q, want %q", branch.Type, "if")
	}
	if len(entry.B["0"]) != 2 {
		t.Errorf("B[\"0\"] length = %d, want 2", len(entry.B["0"]))
	}
	if entry.B["0"][0] != 1 || entry.B["0"][1] != 0 {
		t.Errorf("B[\"0\"] = %v, want [1, 0]", entry.B["0"])
	}
}

func TestParseCoverageFinalJSONFileNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := parseCoverageFinalJSON(dir)
	if err == nil {
		t.Fatal("expected error when coverage-final.json does not exist")
	}
}

func TestParseCoverageFinalJSONInvalid(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "coverage-final.json"), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := parseCoverageFinalJSON(dir)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ---------------------------------------------------------------------------
// collectTSRanges tests
// ---------------------------------------------------------------------------

func TestCollectTSRanges(t *testing.T) {
	fn := &model.Function{
		Name:      "CreateUser",
		File:      "src/handlers/user.ts",
		StartLine: 10,
		EndLine:   30,
	}

	ranges := collectTSRanges(fn)

	if len(ranges) != 1 {
		t.Fatalf("got %d ranges, want 1", len(ranges))
	}

	if ranges[0].file != "src/handlers/user.ts" || ranges[0].funcName != "CreateUser" {
		t.Errorf("ranges[0] = {file: %q, func: %q}, want function", ranges[0].file, ranges[0].funcName)
	}
	if ranges[0].startLine != 10 || ranges[0].endLine != 30 {
		t.Errorf("ranges[0] lines = %d-%d, want 10-30", ranges[0].startLine, ranges[0].endLine)
	}
}

func TestCollectTSRangesEmptyFile(t *testing.T) {
	fn := &model.Function{
		Name: "EmptyFunc",
	}

	ranges := collectTSRanges(fn)
	if len(ranges) != 0 {
		t.Errorf("expected empty ranges for function without file, got %d", len(ranges))
	}
}

// ---------------------------------------------------------------------------
// findCoverageEntry tests
// ---------------------------------------------------------------------------

func TestFindCoverageEntryDirectMatch(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"src/handler.ts": {S: map[string]int{"0": 1}},
	}

	entry := findCoverageEntry("src/handler.ts", data, "/project")
	if entry == nil {
		t.Fatal("expected to find entry for direct match")
	}
	if entry.S["0"] != 1 {
		t.Errorf("S[\"0\"] = %d, want 1", entry.S["0"])
	}
}

func TestFindCoverageEntrySuffixMatch(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"/absolute/path/to/src/handler.ts": {S: map[string]int{"0": 2}},
	}

	entry := findCoverageEntry("src/handler.ts", data, "/project")
	if entry == nil {
		t.Fatal("expected to find entry via suffix match")
	}
	if entry.S["0"] != 2 {
		t.Errorf("S[\"0\"] = %d, want 2", entry.S["0"])
	}
}

func TestFindCoverageEntryAbsProjectMatch(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"/project/src/handler.ts": {S: map[string]int{"0": 3}},
	}

	entry := findCoverageEntry("src/handler.ts", data, "/project")
	if entry == nil {
		t.Fatal("expected to find entry via project root match")
	}
}

func TestFindCoverageEntryNoMatch(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"src/other.ts": {S: map[string]int{"0": 1}},
	}

	entry := findCoverageEntry("src/handler.ts", data, "/project")
	if entry != nil {
		t.Error("expected nil for non-matching file")
	}
}

// ---------------------------------------------------------------------------
// computeTSFuncCoverage tests
// ---------------------------------------------------------------------------

func TestComputeTSFuncCoverage(t *testing.T) {
	data := map[string]coverageFinalEntry{
		"src/handler.ts": {
			StatementMap: map[string]coverageRange{
				"0": {Start: coveragePosition{Line: 10, Column: 0}, End: coveragePosition{Line: 10, Column: 20}},
				"1": {Start: coveragePosition{Line: 15, Column: 0}, End: coveragePosition{Line: 15, Column: 20}},
				"2": {Start: coveragePosition{Line: 50, Column: 0}, End: coveragePosition{Line: 50, Column: 20}}, // outside range
			},
			S: map[string]int{"0": 1, "1": 0, "2": 1},
			BranchMap: map[string]coverageBranch{
				"0": {
					Loc: coverageRange{Start: coveragePosition{Line: 12, Column: 0}, End: coveragePosition{Line: 12, Column: 30}},
					Locations: []coverageRange{
						{Start: coveragePosition{Line: 12, Column: 0}, End: coveragePosition{Line: 12, Column: 15}},
						{Start: coveragePosition{Line: 12, Column: 16}, End: coveragePosition{Line: 12, Column: 30}},
					},
				},
			},
			B:     map[string][]int{"0": {1, 0}},
			FnMap: map[string]coverageFunction{},
			F:     map[string]int{},
		},
	}

	r := tsFuncRange{file: "src/handler.ts", startLine: 10, endLine: 30, funcName: "handler"}
	fc := computeTSFuncCoverage(r, data, "/project")

	// 2 statements in range + 2 branch locations in range = 4 total blocks
	if fc.TotalBlocks != 4 {
		t.Errorf("TotalBlocks = %d, want 4", fc.TotalBlocks)
	}
	// stmt0 covered (1) + stmt1 uncovered (0) + branch loc0 covered (1) + branch loc1 uncovered (0) = 2
	if fc.CoveredBlocks != 2 {
		t.Errorf("CoveredBlocks = %d, want 2", fc.CoveredBlocks)
	}
	if fc.CoveredPct != 50.0 {
		t.Errorf("CoveredPct = %f, want 50.0", fc.CoveredPct)
	}
	if len(fc.UncoveredLines) != 2 {
		t.Errorf("UncoveredLines length = %d, want 2", len(fc.UncoveredLines))
	}
}

func TestComputeTSFuncCoverageNoData(t *testing.T) {
	// Empty coverage data - file not found in coverage, should return 100%.
	data := map[string]coverageFinalEntry{}
	r := tsFuncRange{file: "src/missing.ts", startLine: 1, endLine: 10, funcName: "missing"}
	fc := computeTSFuncCoverage(r, data, "/project")

	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100 (no data)", fc.CoveredPct)
	}
}

func TestComputeTSFuncCoverageEmptyFunc(t *testing.T) {
	// File exists in coverage but function range has no statements.
	data := map[string]coverageFinalEntry{
		"src/handler.ts": {
			StatementMap: map[string]coverageRange{
				"0": {Start: coveragePosition{Line: 50, Column: 0}, End: coveragePosition{Line: 50, Column: 20}},
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
		t.Errorf("CoveredPct = %f, want 100 (empty function)", fc.CoveredPct)
	}
}

// ---------------------------------------------------------------------------
// detectTSCoverageFramework tests
// ---------------------------------------------------------------------------

func writeCoverPkgJSON(t *testing.T, dir string, devDeps map[string]string) {
	t.Helper()
	pkg := struct {
		DevDependencies map[string]string `json:"devDependencies,omitempty"`
	}{DevDependencies: devDeps}
	data, err := json.Marshal(pkg)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDetectTSCoverageFrameworkVitest(t *testing.T) {
	dir := t.TempDir()
	writeCoverPkgJSON(t, dir, map[string]string{"vitest": "^1.0.0"})

	got := detectTSCoverageFramework(dir)
	if got != coverVitest {
		t.Errorf("got %q, want %q", got, coverVitest)
	}
}

func TestDetectTSCoverageFrameworkJest(t *testing.T) {
	dir := t.TempDir()
	writeCoverPkgJSON(t, dir, map[string]string{"jest": "^29.0.0"})

	got := detectTSCoverageFramework(dir)
	if got != coverJest {
		t.Errorf("got %q, want %q", got, coverJest)
	}
}

func TestDetectTSCoverageFrameworkFallback(t *testing.T) {
	dir := t.TempDir()
	// No package.json.
	got := detectTSCoverageFramework(dir)
	if got != coverVitest {
		t.Errorf("got %q, want %q (fallback)", got, coverVitest)
	}
}

func TestDetectTSCoverageFrameworkNoDeps(t *testing.T) {
	dir := t.TempDir()
	writeCoverPkgJSON(t, dir, map[string]string{"typescript": "^5.0.0"})

	got := detectTSCoverageFramework(dir)
	if got != coverVitest {
		t.Errorf("got %q, want %q (no test framework → fallback)", got, coverVitest)
	}
}

// ---------------------------------------------------------------------------
// buildCoverageArgs tests
// ---------------------------------------------------------------------------

func TestBuildCoverageArgsVitest(t *testing.T) {
	args := buildCoverageArgs(coverVitest, "src/handler.test.ts", "/out/coverage")
	expected := []string{
		"vitest", "run", "src/handler.test.ts",
		"--coverage",
		"--coverage.provider=v8",
		"--coverage.reporter=json",
		"--coverage.reportsDirectory=/out/coverage",
	}
	assertCoverArgs(t, args, expected)
}

func TestBuildCoverageArgsJest(t *testing.T) {
	args := buildCoverageArgs(coverJest, "src/handler.test.ts", "/out/coverage")
	expected := []string{
		"jest", "src/handler.test.ts",
		"--coverage",
		"--coverageReporters=json",
		"--coverageDirectory=/out/coverage",
	}
	assertCoverArgs(t, args, expected)
}

func TestBuildCoverageArgsDefault(t *testing.T) {
	// Unknown framework uses vitest args.
	args := buildCoverageArgs("unknown", "src/handler.test.ts", "/out/coverage")
	expected := []string{
		"vitest", "run", "src/handler.test.ts",
		"--coverage",
		"--coverage.provider=v8",
		"--coverage.reporter=json",
		"--coverage.reportsDirectory=/out/coverage",
	}
	assertCoverArgs(t, args, expected)
}

func assertCoverArgs(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %v (len %d), want %v (len %d)", got, len(got), want, len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("arg[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
