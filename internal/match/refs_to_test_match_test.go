package match

import "testing"

func TestRefsToTestMatch_dedups(t *testing.T) {
	refs := []testRef{
		{File: "a_test.go", TestFunc: "TestX"},
		{File: "a_test.go", TestFunc: "TestY"},
		{File: "b_test.go", TestFunc: "TestX"},
	}
	tm, ok := refsToTestMatch(refs)
	if !ok {
		t.Fatal("expected found=true")
	}
	if len(tm.Files) != 2 || tm.Files[0] != "a_test.go" || tm.Files[1] != "b_test.go" {
		t.Errorf("Files = %v, want [a_test.go b_test.go]", tm.Files)
	}
	if len(tm.TestFuncs) != 2 || tm.TestFuncs[0] != "TestX" || tm.TestFuncs[1] != "TestY" {
		t.Errorf("TestFuncs = %v, want [TestX TestY]", tm.TestFuncs)
	}
}

func TestRefsToTestMatch_empty(t *testing.T) {
	if _, ok := refsToTestMatch(nil); ok {
		t.Error("expected found=false for empty refs")
	}
}
