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

// TestMatchTSTopLevelFuncVarArrow covers the `var` arrow-function branch (the
// const|let|var alternative with an arrow body), confirming it is kept.
func TestMatchTSTopLevelFuncVarArrow(t *testing.T) {
	fn, ok := matchTSTopLevelFunc("var doThing = () => {", "pkg", "pkg/x.ts", 7)
	if !ok {
		t.Fatal("expected match for var arrow function")
	}
	if fn.Name != "doThing" {
		t.Errorf("Name = %q, want %q", fn.Name, "doThing")
	}
	if fn.Exported {
		t.Error("expected Exported=false for non-exported var")
	}
	if fn.QualifiedName != "pkg.doThing" {
		t.Errorf("QualifiedName = %q, want pkg.doThing", fn.QualifiedName)
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

// TestMatchTSTopLevelFuncConstFunctionExpr covers the const-with-function-keyword
// path: `const x = function() {}` matches the second alternative (m[2]="run")
// and is kept because the line contains "function".
func TestMatchTSTopLevelFuncConstFunctionExpr(t *testing.T) {
	fn, ok := matchTSTopLevelFunc("const run = function () {", "", "run.ts", 3)
	if !ok {
		t.Fatal("expected match for const function-expression assignment")
	}
	if fn.Name != "run" {
		t.Errorf("Name = %q, want %q", fn.Name, "run")
	}
	if fn.StartLine != 3 || fn.EndLine != 3 {
		t.Errorf("lines = %d-%d, want 3-3", fn.StartLine, fn.EndLine)
	}
}

// TestMatchTSTopLevelFuncDefaultFunction covers the `export default function`
// declaration form (first alternative with the default modifier).
func TestMatchTSTopLevelFuncDefaultFunction(t *testing.T) {
	fn, ok := matchTSTopLevelFunc("export default function main() {", "src", "src/main.ts", 1)
	if !ok {
		t.Fatal("expected match for export default function")
	}
	if fn.Name != "main" {
		t.Errorf("Name = %q, want %q", fn.Name, "main")
	}
	if !fn.Exported {
		t.Error("expected Exported=true for export default function")
	}
}
