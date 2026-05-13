package coverage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestParseCoverLine(t *testing.T) {
	line := "github.com/example/pkg/handler.go:10.2,20.5 3 1"
	b, err := parseCoverLine(line)
	if err != nil {
		t.Fatalf("parseCoverLine: %v", err)
	}
	if b.file != "github.com/example/pkg/handler.go" {
		t.Errorf("file = %q", b.file)
	}
	if b.startLine != 10 {
		t.Errorf("startLine = %d, want 10", b.startLine)
	}
	if b.startCol != 2 {
		t.Errorf("startCol = %d, want 2", b.startCol)
	}
	if b.endLine != 20 {
		t.Errorf("endLine = %d, want 20", b.endLine)
	}
	if b.endCol != 5 {
		t.Errorf("endCol = %d, want 5", b.endCol)
	}
	if b.stmts != 3 {
		t.Errorf("stmts = %d, want 3", b.stmts)
	}
	if b.count != 1 {
		t.Errorf("count = %d, want 1", b.count)
	}
}

func TestParseCoverProfile(t *testing.T) {
	dir := t.TempDir()
	content := `mode: set
github.com/example/pkg/handler.go:10.2,20.5 3 1
github.com/example/pkg/handler.go:22.2,30.5 2 0
github.com/example/pkg/service.go:5.2,15.5 4 1
`
	path := filepath.Join(dir, "cover.out")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	blocks, err := parseCoverProfile(path)
	if err != nil {
		t.Fatalf("parseCoverProfile: %v", err)
	}
	if len(blocks) != 3 {
		t.Fatalf("got %d blocks, want 3", len(blocks))
	}
	if blocks[1].count != 0 {
		t.Error("second block should be uncovered")
	}
}

// ---------------------------------------------------------------------------
// Factory tests
// ---------------------------------------------------------------------------

func TestNewCheckerGo(t *testing.T) {
	c := NewChecker("go")
	if _, ok := c.(*GoChecker); !ok {
		t.Fatalf("NewChecker(\"go\") returned %T, want *GoChecker", c)
	}
}

func TestNewCheckerTS(t *testing.T) {
	c := NewChecker("typescript")
	if _, ok := c.(*TSChecker); !ok {
		t.Fatalf("NewChecker(\"typescript\") returned %T, want *TSChecker", c)
	}
}

func TestNewCheckerPython(t *testing.T) {
	c := NewChecker("python")
	if _, ok := c.(*PyChecker); !ok {
		t.Fatalf("NewChecker(\"python\") returned %T, want *PyChecker", c)
	}
}

func TestNewCheckerUnsupported(t *testing.T) {
	c := NewChecker("rust")
	u, ok := c.(*UnsupportedChecker)
	if !ok {
		t.Fatalf("NewChecker(\"rust\") returned %T, want *UnsupportedChecker", c)
	}
	if u.Lang != "rust" {
		t.Errorf("UnsupportedChecker.Lang = %q, want %q", u.Lang, "rust")
	}

	_, err := c.Check("/tmp", "main_test.rs", &model.Function{})
	if err == nil {
		t.Fatal("UnsupportedChecker.Check should return an error")
	}
	var unsupErr *ErrUnsupported
	if e, ok := err.(*ErrUnsupported); !ok {
		t.Errorf("error type = %T, want *ErrUnsupported", err)
	} else {
		unsupErr = e
	}
	if unsupErr != nil && !strings.Contains(unsupErr.Error(), "rust") {
		t.Errorf("error should mention 'rust': %v", unsupErr)
	}
}

// ---------------------------------------------------------------------------
// collectRanges tests (Go checker)
// ---------------------------------------------------------------------------

func TestCollectRanges(t *testing.T) {
	fn := &model.Function{
		Name:      "GetUser",
		File:      "internal/handler/user.go",
		StartLine: 10,
		EndLine:   30,
	}

	ranges := collectRanges(fn)

	if len(ranges) != 1 {
		t.Fatalf("got %d ranges, want 1", len(ranges))
	}

	if ranges[0].file != "internal/handler/user.go" || ranges[0].funcName != "GetUser" {
		t.Errorf("ranges[0] = {file: %q, func: %q}, want function", ranges[0].file, ranges[0].funcName)
	}
	if ranges[0].startLine != 10 || ranges[0].endLine != 30 {
		t.Errorf("ranges[0] lines = %d-%d, want 10-30", ranges[0].startLine, ranges[0].endLine)
	}
}

func TestCollectRangesEmptyFile(t *testing.T) {
	fn := &model.Function{
		Name: "Empty",
	}

	ranges := collectRanges(fn)
	if len(ranges) != 0 {
		t.Errorf("expected empty ranges, got %d", len(ranges))
	}
}

// ---------------------------------------------------------------------------
// computeFuncCoverage tests
// ---------------------------------------------------------------------------

func TestComputeFuncCoverage(t *testing.T) {
	blocks := []coverBlock{
		{file: "github.com/example/pkg/handler.go", startLine: 10, endLine: 15, count: 1},
		{file: "github.com/example/pkg/handler.go", startLine: 16, endLine: 20, count: 0},
		{file: "github.com/example/pkg/handler.go", startLine: 22, endLine: 25, count: 1},
		{file: "github.com/example/pkg/handler.go", startLine: 50, endLine: 55, count: 1}, // outside range
		{file: "github.com/example/pkg/service.go", startLine: 10, endLine: 20, count: 1}, // different file
	}

	r := funcRange{file: "pkg/handler.go", startLine: 10, endLine: 25, funcName: "HandleRequest"}
	fc := computeFuncCoverage(r, blocks, "/project")

	if fc.TotalBlocks != 3 {
		t.Errorf("TotalBlocks = %d, want 3", fc.TotalBlocks)
	}
	if fc.CoveredBlocks != 2 {
		t.Errorf("CoveredBlocks = %d, want 2", fc.CoveredBlocks)
	}
	expectedPct := float64(2) / float64(3) * 100
	if fc.CoveredPct != expectedPct {
		t.Errorf("CoveredPct = %f, want %f", fc.CoveredPct, expectedPct)
	}
	if len(fc.UncoveredLines) != 1 || fc.UncoveredLines[0] != 16 {
		t.Errorf("UncoveredLines = %v, want [16]", fc.UncoveredLines)
	}
}

func TestComputeFuncCoverageNoBlocks(t *testing.T) {
	// No matching blocks — treated as 100% (empty function).
	r := funcRange{file: "pkg/empty.go", startLine: 1, endLine: 5, funcName: "Empty"}
	fc := computeFuncCoverage(r, nil, "/project")
	if fc.CoveredPct != 100 {
		t.Errorf("CoveredPct = %f, want 100 (no blocks)", fc.CoveredPct)
	}
}

func TestOverlaps(t *testing.T) {
	tests := []struct {
		blockFile  string
		targetFile string
		blockStart int
		blockEnd   int
		funcStart  int
		funcEnd    int
		want       bool
	}{
		{"github.com/example/pkg/handler.go", "pkg/handler.go", 15, 20, 10, 25, true},
		{"github.com/example/pkg/handler.go", "pkg/handler.go", 5, 8, 10, 25, false},
		{"github.com/example/pkg/handler.go", "pkg/service.go", 15, 20, 10, 25, false},
		{"github.com/example/pkg/handler.go", "pkg/handler.go", 10, 12, 10, 25, true},
		{"github.com/example/pkg/handler.go", "pkg/handler.go", 25, 30, 10, 25, true},
		{"github.com/example/pkg/handler.go", "pkg/handler.go", 26, 30, 10, 25, false},
	}

	for i, tt := range tests {
		got := overlaps(tt.blockFile, tt.targetFile, tt.blockStart, tt.blockEnd, tt.funcStart, tt.funcEnd)
		if got != tt.want {
			t.Errorf("case %d: overlaps(%q, %q, %d, %d, %d, %d) = %v, want %v",
				i, tt.blockFile, tt.targetFile, tt.blockStart, tt.blockEnd, tt.funcStart, tt.funcEnd, got, tt.want)
		}
	}
}
