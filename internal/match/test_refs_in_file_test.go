package match

import "testing"

func TestTestRefsInFile(t *testing.T) {
	path := writeTmpFile(t, "x_test.go", `package p
import "testing"
func TestA(t *testing.T) { Alpha() }
func TestB(t *testing.T) { Beta() }
func helper() { NotATest() }
`)
	funcs, err := parseTestFileFuncs(path)
	if err != nil {
		t.Fatal(err)
	}
	result := testRefsInFile(funcs)

	// Only TestA and TestB are keyed (helper is not a Test entrypoint).
	if len(result) != 2 {
		t.Fatalf("got %d test entries, want 2: %v", len(result), result)
	}
	if _, ok := result["helper"]; ok {
		t.Errorf("non-test helper must not be a key: %v", result)
	}
	if !hasAll(refNames(result["TestA"]), "Alpha") {
		t.Errorf("TestA refs = %v, want Alpha", result["TestA"])
	}
	if !hasAll(refNames(result["TestB"]), "Beta") {
		t.Errorf("TestB refs = %v, want Beta", result["TestB"])
	}
}

func TestTestRefsInFileNoTests(t *testing.T) {
	path := writeTmpFile(t, "notest_test.go", `package p
func helper() { Foo() }
`)
	funcs, err := parseTestFileFuncs(path)
	if err != nil {
		t.Fatal(err)
	}
	result := testRefsInFile(funcs)
	if len(result) != 0 {
		t.Errorf("got %v, want empty (no Test* funcs)", result)
	}
}
