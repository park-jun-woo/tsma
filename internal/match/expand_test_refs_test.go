package match

import "testing"

func TestExpandTestRefsDirectAndHelper(t *testing.T) {
	path := writeTmpFile(t, "x_test.go", `package p
import "testing"
func TestX(t *testing.T) {
	Direct()
	helper(t)
}
func helper(t *testing.T) { ViaHelper() }
`)
	funcs, err := parseTestFileFuncs(path)
	if err != nil {
		t.Fatal(err)
	}
	refs := expandTestRefs(funcs["TestX"], funcs)

	// Direct callee, the helper itself, and the helper's 1-hop callee.
	if !hasAll(refs, "Direct", "helper", "ViaHelper") {
		t.Errorf("expandTestRefs = %v, want Direct,helper,ViaHelper", refs)
	}
}

func TestExpandTestRefsNoDeepIndirection(t *testing.T) {
	path := writeTmpFile(t, "deep_test.go", `package p
import "testing"
func TestDeep(t *testing.T) { hopOne(t) }
func hopOne(t *testing.T) { hopTwo(t) }
func hopTwo(t *testing.T) { DeepTarget() }
`)
	funcs, err := parseTestFileFuncs(path)
	if err != nil {
		t.Fatal(err)
	}
	refs := expandTestRefs(funcs["TestDeep"], funcs)

	if !hasAll(refs, "hopOne", "hopTwo") {
		t.Errorf("expandTestRefs = %v, want hopOne,hopTwo (1-hop)", refs)
	}
	if _, ok := refs["DeepTarget"]; ok {
		t.Errorf("DeepTarget should NOT be reached past 1-hop: %v", refs)
	}
}

func TestExpandTestRefsDoesNotExpandTestHelper(t *testing.T) {
	// A directly-called func that is itself a Test* must not be expanded.
	path := writeTmpFile(t, "th_test.go", `package p
import "testing"
func TestOuter(t *testing.T) { TestInner(t) }
func TestInner(t *testing.T) { Hidden() }
`)
	funcs, err := parseTestFileFuncs(path)
	if err != nil {
		t.Fatal(err)
	}
	refs := expandTestRefs(funcs["TestOuter"], funcs)
	if _, ok := refs["TestInner"]; !ok {
		t.Errorf("TestInner should be a direct ref: %v", refs)
	}
	if _, ok := refs["Hidden"]; ok {
		t.Errorf("Hidden must not be expanded through a Test* func: %v", refs)
	}
}
