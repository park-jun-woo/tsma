package coverage

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestCollectCsRanges(t *testing.T) {
	fn := &model.Function{File: "App/Calculator.cs", Name: "Add", StartLine: 5, EndLine: 7}
	ranges := collectCsRanges(fn)
	if len(ranges) != 1 {
		t.Fatalf("got %d ranges, want 1", len(ranges))
	}
	r := ranges[0]
	if r.file != "App/Calculator.cs" || r.startLine != 5 || r.endLine != 7 || r.funcName != "Add" {
		t.Errorf("range = %+v", r)
	}
}

func TestCollectCsRangesNoFile(t *testing.T) {
	if r := collectCsRanges(&model.Function{}); r != nil {
		t.Errorf("expected nil for empty file, got %+v", r)
	}
}
