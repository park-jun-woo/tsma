package model

import "testing"

func TestFunctionStructFields(t *testing.T) {
	f := Function{
		QualifiedName: "pkg.Login",
		Name:          "Login",
		File:          "handler.go",
		StartLine:     10,
		EndLine:       30,
		Exported:      true,
		TestFile:      "handler_test.go",
		TestMtime:     "2025-01-01T00:00:00Z",
		Status:        StatusTodo,
		CoveragePct:   85.5,
		Attempt:       2,
		FailOutput:    "assertion failed",
	}

	if f.QualifiedName != "pkg.Login" {
		t.Errorf("QualifiedName = %q, want %q", f.QualifiedName, "pkg.Login")
	}
	if f.Name != "Login" {
		t.Errorf("Name = %q, want %q", f.Name, "Login")
	}
	if f.File != "handler.go" {
		t.Errorf("File = %q, want %q", f.File, "handler.go")
	}
	if f.StartLine != 10 {
		t.Errorf("StartLine = %d, want 10", f.StartLine)
	}
	if f.EndLine != 30 {
		t.Errorf("EndLine = %d, want 30", f.EndLine)
	}
	if !f.Exported {
		t.Error("Exported = false, want true")
	}
	if f.TestFile != "handler_test.go" {
		t.Errorf("TestFile = %q, want %q", f.TestFile, "handler_test.go")
	}
	if f.TestMtime != "2025-01-01T00:00:00Z" {
		t.Errorf("TestMtime = %q, want %q", f.TestMtime, "2025-01-01T00:00:00Z")
	}
	if f.Status != StatusTodo {
		t.Errorf("Status = %q, want %q", f.Status, StatusTodo)
	}
	if f.CoveragePct != 85.5 {
		t.Errorf("CoveragePct = %f, want 85.5", f.CoveragePct)
	}
	if f.Attempt != 2 {
		t.Errorf("Attempt = %d, want 2", f.Attempt)
	}
	if f.FailOutput != "assertion failed" {
		t.Errorf("FailOutput = %q, want %q", f.FailOutput, "assertion failed")
	}
}

func TestFunctionZeroValue(t *testing.T) {
	var f Function
	if f.QualifiedName != "" {
		t.Errorf("zero QualifiedName = %q, want empty", f.QualifiedName)
	}
	if f.Exported {
		t.Error("zero Exported = true, want false")
	}
	if f.Status != "" {
		t.Errorf("zero Status = %q, want empty", f.Status)
	}
	if f.CoveragePct != 0 {
		t.Errorf("zero CoveragePct = %f, want 0", f.CoveragePct)
	}
	if f.Attempt != 0 {
		t.Errorf("zero Attempt = %d, want 0", f.Attempt)
	}
}
