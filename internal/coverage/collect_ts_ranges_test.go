package coverage

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestCollectTSRangesWithFile(t *testing.T) {
	fn := &model.Function{
		Name:      "createUser",
		File:      "src/handlers/user.ts",
		StartLine: 10,
		EndLine:   30,
	}

	ranges := collectTSRanges(fn)
	if len(ranges) != 1 {
		t.Fatalf("got %d ranges, want 1", len(ranges))
	}
	if ranges[0].file != fn.File {
		t.Errorf("file = %q, want %q", ranges[0].file, fn.File)
	}
	if ranges[0].funcName != fn.Name {
		t.Errorf("funcName = %q, want %q", ranges[0].funcName, fn.Name)
	}
	if ranges[0].startLine != 10 || ranges[0].endLine != 30 {
		t.Errorf("lines = %d-%d, want 10-30", ranges[0].startLine, ranges[0].endLine)
	}
}

func TestCollectTSRangesWithoutFile(t *testing.T) {
	fn := &model.Function{Name: "noFile"}
	ranges := collectTSRanges(fn)
	if len(ranges) != 0 {
		t.Errorf("expected empty ranges, got %d", len(ranges))
	}
}
