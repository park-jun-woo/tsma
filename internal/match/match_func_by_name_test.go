package match

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

func TestMatchFuncByNameHit(t *testing.T) {
	idx := &PkgTestIndex{refs: map[string][]testRef{
		"Generate": {{File: "g_test.go", TestFunc: "TestGen"}},
	}}
	fn := &model.Function{Name: "Generate"}
	refs, ok := MatchFuncByName(idx, fn)
	if !ok || len(refs) != 1 || refs[0].TestFunc != "TestGen" {
		t.Fatalf("MatchFuncByName = %v ok=%v, want [TestGen]", refs, ok)
	}
}

func TestMatchFuncByNameMiss(t *testing.T) {
	idx := &PkgTestIndex{refs: map[string][]testRef{"Generate": {{TestFunc: "TestGen"}}}}
	fn := &model.Function{Name: "Nope"}
	if refs, ok := MatchFuncByName(idx, fn); ok || refs != nil {
		t.Errorf("MatchFuncByName(Nope) = %v ok=%v, want nil,false", refs, ok)
	}
}

func TestMatchFuncByNameNilIndex(t *testing.T) {
	if refs, ok := MatchFuncByName(nil, &model.Function{Name: "X"}); ok || refs != nil {
		t.Errorf("nil idx = %v ok=%v, want nil,false", refs, ok)
	}
}

func TestMatchFuncByNameNilFunc(t *testing.T) {
	idx := &PkgTestIndex{refs: map[string][]testRef{"X": {{TestFunc: "TestX"}}}}
	if refs, ok := MatchFuncByName(idx, nil); ok || refs != nil {
		t.Errorf("nil fn = %v ok=%v, want nil,false", refs, ok)
	}
}
