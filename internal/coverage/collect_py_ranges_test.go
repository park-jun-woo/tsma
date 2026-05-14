package coverage

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestCollectPyRangesWithFile(t *testing.T) {
	fn := &model.Function{
		Name:      "create_order",
		File:      "handlers/order.py",
		StartLine: 5,
		EndLine:   25,
	}

	ranges := collectPyRanges(fn)
	if len(ranges) != 1 {
		t.Fatalf("got %d ranges, want 1", len(ranges))
	}
	if ranges[0].file != fn.File {
		t.Errorf("file = %q, want %q", ranges[0].file, fn.File)
	}
	if ranges[0].funcName != fn.Name {
		t.Errorf("funcName = %q, want %q", ranges[0].funcName, fn.Name)
	}
	if ranges[0].startLine != 5 || ranges[0].endLine != 25 {
		t.Errorf("lines = %d-%d, want 5-25", ranges[0].startLine, ranges[0].endLine)
	}
}

func TestCollectPyRangesWithoutFile(t *testing.T) {
	fn := &model.Function{Name: "empty"}
	ranges := collectPyRanges(fn)
	if len(ranges) != 0 {
		t.Errorf("expected empty ranges, got %d", len(ranges))
	}
}
