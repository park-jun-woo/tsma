package cli

import (
	"testing"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
)

// TestSetTestFiles_recordsFilesAndRepresentative verifies setTestFiles copies
// the full file set into TestFiles and uses the first file as the representative
// TestFile.
func TestSetTestFiles_recordsFilesAndRepresentative(t *testing.T) {
	fn := &model.Function{Name: "Foo"}
	tm := match.TestMatch{Files: []string{"a_test.go", "b_test.go"}}

	setTestFiles(fn, tm)

	if len(fn.TestFiles) != 2 || fn.TestFiles[0] != "a_test.go" || fn.TestFiles[1] != "b_test.go" {
		t.Errorf("TestFiles = %v, want [a_test.go b_test.go]", fn.TestFiles)
	}
	if fn.TestFile != "a_test.go" {
		t.Errorf("TestFile = %q, want %q", fn.TestFile, "a_test.go")
	}
}

// TestSetTestFiles_emptyIsNoOp verifies a match with no files leaves the
// function untouched.
func TestSetTestFiles_emptyIsNoOp(t *testing.T) {
	fn := &model.Function{Name: "Foo", TestFile: "old_test.go", TestFiles: []string{"old_test.go"}}

	setTestFiles(fn, match.TestMatch{})

	if fn.TestFile != "old_test.go" {
		t.Errorf("TestFile = %q, want unchanged old_test.go", fn.TestFile)
	}
	if len(fn.TestFiles) != 1 || fn.TestFiles[0] != "old_test.go" {
		t.Errorf("TestFiles = %v, want unchanged [old_test.go]", fn.TestFiles)
	}
}
