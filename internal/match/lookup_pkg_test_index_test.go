package match

import "testing"

func TestLookupHit(t *testing.T) {
	idx := &PkgTestIndex{refs: map[string][]testRef{
		"Foo": {{File: "x_test.go", TestFunc: "TestFoo"}},
	}}
	refs, ok := idx.Lookup("Foo")
	if !ok || len(refs) != 1 || refs[0].TestFunc != "TestFoo" {
		t.Fatalf("Lookup(Foo) = %v ok=%v, want [TestFoo]", refs, ok)
	}
}

func TestLookupMiss(t *testing.T) {
	idx := &PkgTestIndex{refs: map[string][]testRef{"Foo": {{TestFunc: "TestFoo"}}}}
	if refs, ok := idx.Lookup("Bar"); ok || refs != nil {
		t.Errorf("Lookup(Bar) = %v ok=%v, want nil,false", refs, ok)
	}
}

func TestLookupEmptySlice(t *testing.T) {
	idx := &PkgTestIndex{refs: map[string][]testRef{"Foo": {}}}
	if refs, ok := idx.Lookup("Foo"); ok || refs != nil {
		t.Errorf("Lookup(Foo) with empty slice = %v ok=%v, want nil,false", refs, ok)
	}
}

func TestLookupNilReceiver(t *testing.T) {
	var idx *PkgTestIndex
	if refs, ok := idx.Lookup("Foo"); ok || refs != nil {
		t.Errorf("nil index Lookup = %v ok=%v, want nil,false", refs, ok)
	}
}

func TestLookupNilRefsMap(t *testing.T) {
	idx := &PkgTestIndex{}
	if refs, ok := idx.Lookup("Foo"); ok || refs != nil {
		t.Errorf("nil refs map Lookup = %v ok=%v, want nil,false", refs, ok)
	}
}
