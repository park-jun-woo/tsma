package coverage

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestCollectRsRanges(t *testing.T) {
	fn := &model.Function{File: "src/lib.rs", Name: "classify", StartLine: 1, EndLine: 5}
	ranges := collectRsRanges(fn)
	if len(ranges) != 1 {
		t.Fatalf("got %d ranges, want 1", len(ranges))
	}
	r := ranges[0]
	if r.file != "src/lib.rs" || r.startLine != 1 || r.endLine != 5 || r.funcName != "classify" {
		t.Errorf("range = %+v", r)
	}
}

func TestCollectRsRangesNoFile(t *testing.T) {
	if r := collectRsRanges(&model.Function{}); r != nil {
		t.Errorf("expected nil for empty file, got %+v", r)
	}
}
