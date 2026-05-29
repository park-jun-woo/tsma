package coverage

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestCollectJavaRanges(t *testing.T) {
	fn := &model.Function{File: "src/main/java/p/Foo.java", Name: "add", StartLine: 5, EndLine: 7}
	ranges := collectJavaRanges(fn)
	if len(ranges) != 1 {
		t.Fatalf("got %d ranges, want 1", len(ranges))
	}
	r := ranges[0]
	if r.file != "src/main/java/p/Foo.java" || r.startLine != 5 || r.endLine != 7 || r.funcName != "add" {
		t.Errorf("range = %+v", r)
	}
}

func TestCollectJavaRangesNoFile(t *testing.T) {
	if r := collectJavaRanges(&model.Function{}); r != nil {
		t.Errorf("expected nil for empty file, got %+v", r)
	}
}
