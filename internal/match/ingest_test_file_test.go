package match

import "testing"

func TestIngestTestFileRecordsRefs(t *testing.T) {
	abs := writeTmpFile(t, "x_test.go", `package p
import "testing"
func TestX(t *testing.T) { Foo() }
`)
	idx := &PkgTestIndex{refs: make(map[string][]testRef)}
	ingestTestFile(idx, abs, "rel/x_test.go")

	refs, ok := idx.Lookup("Foo")
	if !ok || len(refs) != 1 {
		t.Fatalf("Lookup(Foo) = %v ok=%v, want 1 ref", refs, ok)
	}
	if refs[0].TestFunc != "TestX" {
		t.Errorf("TestFunc = %q, want TestX", refs[0].TestFunc)
	}
	if refs[0].File != "rel/x_test.go" {
		t.Errorf("File = %q, want rel/x_test.go (the rel arg)", refs[0].File)
	}
}

func TestIngestTestFileSkipsUnparseable(t *testing.T) {
	abs := writeTmpFile(t, "bad_test.go", `package p
func TestX(t *testing.T) { not valid go
`)
	idx := &PkgTestIndex{refs: make(map[string][]testRef)}
	// Must not panic and must leave the index empty.
	ingestTestFile(idx, abs, "rel/bad_test.go")
	if len(idx.refs) != 0 {
		t.Errorf("unparseable file should add nothing, got %v", idx.refs)
	}
}

func TestIngestTestFileAppends(t *testing.T) {
	a := writeTmpFile(t, "a_test.go", `package p
import "testing"
func TestA(t *testing.T) { Shared() }
`)
	b := writeTmpFile(t, "b_test.go", `package p
import "testing"
func TestB(t *testing.T) { Shared() }
`)
	idx := &PkgTestIndex{refs: make(map[string][]testRef)}
	ingestTestFile(idx, a, "a_test.go")
	ingestTestFile(idx, b, "b_test.go")

	refs, ok := idx.Lookup("Shared")
	if !ok || len(refs) != 2 {
		t.Fatalf("Lookup(Shared) = %v ok=%v, want 2 refs (appended)", refs, ok)
	}
}
