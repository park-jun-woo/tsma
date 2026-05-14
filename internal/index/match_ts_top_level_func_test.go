package index

import "testing"

func TestMatchTSTopLevelFuncExportedFunction(t *testing.T) {
	fn, ok := matchTSTopLevelFunc("export async function handleLogin(req: Request) {", "src/api", "src/api/handler.ts", 5)
	if !ok {
		t.Fatal("expected match for exported async function")
	}
	if fn.Name != "handleLogin" {
		t.Errorf("Name = %q, want %q", fn.Name, "handleLogin")
	}
	if fn.QualifiedName != "src/api.handleLogin" {
		t.Errorf("QualifiedName = %q, want %q", fn.QualifiedName, "src/api.handleLogin")
	}
	if !fn.Exported {
		t.Error("expected Exported=true")
	}
	if fn.StartLine != 5 {
		t.Errorf("StartLine = %d, want 5", fn.StartLine)
	}
}

func TestMatchTSTopLevelFuncNonExported(t *testing.T) {
	fn, ok := matchTSTopLevelFunc("function helperFunc() {", "utils", "utils/helper.ts", 10)
	if !ok {
		t.Fatal("expected match for non-exported function")
	}
	if fn.Name != "helperFunc" {
		t.Errorf("Name = %q, want %q", fn.Name, "helperFunc")
	}
	if fn.Exported {
		t.Error("expected Exported=false for non-export function")
	}
}

func TestMatchTSTopLevelFuncConstArrow(t *testing.T) {
	fn, ok := matchTSTopLevelFunc("export const formatDate = (d: Date) => {", "", "utils.ts", 1)
	if !ok {
		t.Fatal("expected match for const arrow function")
	}
	if fn.Name != "formatDate" {
		t.Errorf("Name = %q, want %q", fn.Name, "formatDate")
	}
	if !fn.Exported {
		t.Error("expected Exported=true for export const")
	}
}

func TestMatchTSTopLevelFuncNoMatch(t *testing.T) {
	_, ok := matchTSTopLevelFunc("return result;", "", "file.ts", 1)
	if ok {
		t.Error("expected no match for non-function line")
	}
}

func TestMatchTSTopLevelFuncConstAssignment(t *testing.T) {
	// const x = 42 should NOT match since Phase007 filters non-function consts
	_, ok := matchTSTopLevelFunc("const x = 42;", "", "file.ts", 1)
	if ok {
		t.Error("expected no match for const non-function assignment")
	}
}

func TestMatchTSTopLevelFuncLetArrow(t *testing.T) {
	fn, ok := matchTSTopLevelFunc("let handler = (req) => {", "", "handler.ts", 1)
	if !ok {
		t.Fatal("expected match for let arrow function")
	}
	if fn.Name != "handler" {
		t.Errorf("Name = %q, want %q", fn.Name, "handler")
	}
}
