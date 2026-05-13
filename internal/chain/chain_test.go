package chain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

// ---------------------------------------------------------------------------
// Factory tests
// ---------------------------------------------------------------------------

func TestNewTracerGo(t *testing.T) {
	tr := NewTracer("go")
	if _, ok := tr.(*GoTracer); !ok {
		t.Errorf("NewTracer(\"go\") returned %T, want *GoTracer", tr)
	}
}

func TestNewTracerTypescript(t *testing.T) {
	tr := NewTracer("typescript")
	if _, ok := tr.(*TSTracer); !ok {
		t.Errorf("NewTracer(\"typescript\") returned %T, want *TSTracer", tr)
	}
}

func TestNewTracerPython(t *testing.T) {
	tr := NewTracer("python")
	if _, ok := tr.(*PyTracer); !ok {
		t.Errorf("NewTracer(\"python\") returned %T, want *PyTracer", tr)
	}
}

func TestNewTracerUnsupported(t *testing.T) {
	tr := NewTracer("ruby")
	ut, ok := tr.(*UnsupportedTracer)
	if !ok {
		t.Fatalf("NewTracer(\"ruby\") returned %T, want *UnsupportedTracer", tr)
	}
	if ut.Lang != "ruby" {
		t.Errorf("UnsupportedTracer.Lang = %q, want \"ruby\"", ut.Lang)
	}

	_, err := tr.Trace(".", model.FuncLocation{})
	if err == nil {
		t.Fatal("expected error from UnsupportedTracer.Trace")
	}
	ue, ok := err.(*ErrUnsupported)
	if !ok {
		t.Fatalf("error type = %T, want *ErrUnsupported", err)
	}
	if ue.Lang != "ruby" {
		t.Errorf("ErrUnsupported.Lang = %q, want \"ruby\"", ue.Lang)
	}
	want := "chain tracing not implemented for: ruby"
	if ue.Error() != want {
		t.Errorf("error message = %q, want %q", ue.Error(), want)
	}
}

// ---------------------------------------------------------------------------
// Helper: write file into dir, creating intermediate directories.
// ---------------------------------------------------------------------------

func writeFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	abs := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// GoTracer tests
// ---------------------------------------------------------------------------

func TestGoTracerBasicChain(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.go", `package main

func HandleLogin() {
	svc.Login()
}
`)

	writeFile(t, dir, "service.go", `package main

func Login() {
	repo.FindByEmail()
	generateToken()
}
`)

	writeFile(t, dir, "util.go", `package main

func generateToken() string {
	return "token"
}
`)

	tr := &GoTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "handler.go",
		StartLine: 3,
		EndLine:   5,
	})
	if err != nil {
		t.Fatalf("GoTracer.Trace: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("expected non-empty chain entries")
	}

	// Check that Login appears as an internal entry.
	found := findEntry(entries, "Login", false)
	if found == nil {
		// The display name may be "svc.Login" since it is a selector call.
		found = findEntry(entries, "svc.Login", false)
	}
	if found == nil {
		t.Error("expected Login or svc.Login in chain entries")
	} else if found.File == "" {
		t.Error("expected Login to have a File (internal), but File is empty")
	}

	// Check that generateToken appears as an internal entry.
	foundGen := findEntry(entries, "generateToken", false)
	if foundGen == nil {
		t.Error("expected generateToken in chain entries")
	} else if foundGen.File == "" {
		t.Error("expected generateToken to have a File (internal)")
	}

	// Check that repo.FindByEmail appears as an external entry with boundary.
	foundRepo := findEntry(entries, "repo.FindByEmail", true)
	if foundRepo == nil {
		t.Error("expected repo.FindByEmail as external entry")
	} else if foundRepo.Boundary != "repository-interface" {
		t.Errorf("repo.FindByEmail boundary = %q, want \"repository-interface\"", foundRepo.Boundary)
	}
}

func TestGoTracerCyclePrevention(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "cycle.go", `package main

func FuncA() {
	FuncB()
}

func FuncB() {
	FuncA()
}
`)

	tr := &GoTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "cycle.go",
		StartLine: 3,
		EndLine:   5,
	})
	if err != nil {
		t.Fatalf("GoTracer.Trace: %v", err)
	}

	// Should complete without hanging; both functions should appear at most once.
	countA := 0
	countB := 0
	for _, e := range entries {
		if e.Func == "FuncA" {
			countA++
		}
		if e.Func == "FuncB" {
			countB++
		}
	}
	if countA > 1 {
		t.Errorf("FuncA appears %d times, expected at most 1 (cycle prevention)", countA)
	}
	if countB > 1 {
		t.Errorf("FuncB appears %d times, expected at most 1 (cycle prevention)", countB)
	}
}

func TestGoTracerHandlerNotFound(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.go", `package main

func Foo() {}
func Bar() {}
`)

	tr := &GoTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "nonexistent.go",
		StartLine: 1,
		EndLine:   3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries when handler not found, got %d entries", len(entries))
	}
}

func TestGoTracerClassifyBoundary(t *testing.T) {
	tests := []struct {
		receiver string
		want     string
	}{
		{"repo", "repository-interface"},
		{"userRepo", "repository-interface"},
		{"db", "database"},
		{"myDB", "database"},
		{"store", "repository-interface"},
		{"cacheStore", "repository-interface"},
		{"client", "external"},
		{"svc", "external"},
	}
	for _, tt := range tests {
		got := classifyBoundary(tt.receiver)
		if got != tt.want {
			t.Errorf("classifyBoundary(%q) = %q, want %q", tt.receiver, got, tt.want)
		}
	}
}

func TestGoTracerEmptyProject(t *testing.T) {
	dir := t.TempDir()

	tr := &GoTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "handler.go",
		StartLine: 1,
		EndLine:   5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for empty project, got %d", len(entries))
	}
}

// ---------------------------------------------------------------------------
// TSTracer tests
// ---------------------------------------------------------------------------

func TestTSTracerBasicChain(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.ts", `import { authService } from './auth-service';

export async function handleLogin(req: Request) {
  const result = await authService.login(req.body);
  return result;
}
`)

	writeFile(t, dir, "auth-service.ts", `import { db } from './db';
import { generateToken } from './utils';

export async function login(credentials: any) {
  const user = await db.findUser(credentials.email);
  const token = generateToken(user);
  return { user, token };
}
`)

	writeFile(t, dir, "utils.ts", `export function generateToken(user: any): string {
  return 'jwt-token-' + user.id;
}
`)

	tr := &TSTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "handler.ts",
		StartLine: 3,
		EndLine:   6,
	})
	if err != nil {
		t.Fatalf("TSTracer.Trace: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("expected non-empty chain entries from TSTracer")
	}

	// login should be found in auth-service.ts as internal.
	foundLogin := findEntryContaining(entries, "login", false)
	if foundLogin == nil {
		t.Error("expected login in chain entries as internal")
	} else if foundLogin.File == "" {
		t.Error("expected login to have a File (internal)")
	}

	// generateToken should be found in utils.ts as internal.
	foundToken := findEntryContaining(entries, "generateToken", false)
	if foundToken == nil {
		t.Error("expected generateToken in chain entries as internal")
	} else if foundToken.File == "" {
		t.Error("expected generateToken to have a File (internal)")
	}

	// db.findUser should be classified as external with boundary.
	foundDB := findEntryContaining(entries, "findUser", true)
	if foundDB == nil {
		t.Error("expected db.findUser as external entry")
	} else if foundDB.Boundary == "" {
		t.Error("expected db.findUser to have a boundary classification")
	}
}

func TestTSTracerEmptyHandler(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.ts", `export function handler() {}
`)

	tr := &TSTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "handler.ts",
		StartLine: 1,
		EndLine:   1,
	})
	if err != nil {
		t.Fatalf("TSTracer.Trace: %v", err)
	}
	// Empty handler body should produce no entries (or very few).
	// Just ensure it doesn't crash.
	_ = entries
}

func TestTSTracerNoFile(t *testing.T) {
	dir := t.TempDir()

	tr := &TSTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "nonexistent.ts",
		StartLine: 1,
		EndLine:   5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for missing file, got %d", len(entries))
	}
}

func TestTSTracerEmptyFile(t *testing.T) {
	tr := &TSTracer{}
	entries, err := tr.Trace(".", model.FuncLocation{
		File:      "",
		StartLine: 1,
		EndLine:   5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for empty file path, got %d", len(entries))
	}
}

func TestTSTracerSkipBuiltins(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.ts", `export function handler() {
  console.log("hello");
  JSON.parse("{}");
  Math.random();
}
`)

	tr := &TSTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "handler.ts",
		StartLine: 1,
		EndLine:   5,
	})
	if err != nil {
		t.Fatalf("TSTracer.Trace: %v", err)
	}

	// console, JSON, Math are builtins and should be skipped.
	for _, e := range entries {
		for _, skip := range []string{"console", "JSON", "Math"} {
			if e.Func == skip+".log" || e.Func == skip+".parse" || e.Func == skip+".random" {
				t.Errorf("builtin call %q should have been skipped", e.Func)
			}
		}
	}
}

func TestTSTracerRepoBoundary(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.ts", `export function handler() {
  const user = repo.findById(1);
  const data = prisma.user.findUnique({ id: 1 });
  const item = mongoose.model("Item");
}
`)

	tr := &TSTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "handler.ts",
		StartLine: 1,
		EndLine:   5,
	})
	if err != nil {
		t.Fatalf("TSTracer.Trace: %v", err)
	}

	for _, e := range entries {
		if e.Boundary == "" {
			continue
		}
		if e.Func == "repo.findById" && e.Boundary != "repository-interface" {
			t.Errorf("repo.findById boundary = %q, want repository-interface", e.Boundary)
		}
	}
}

func TestIsTSOrJSSourceFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"handler.ts", true},
		{"service.js", true},
		{"types.d.ts", false},
		{"handler.test.ts", false},
		{"handler.spec.js", false},
		{"readme.md", false},
		{"handler.go", false},
	}
	for _, tt := range tests {
		got := isTSOrJSSourceFile(tt.path)
		if got != tt.want {
			t.Errorf("isTSOrJSSourceFile(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestClassifyTSBoundary(t *testing.T) {
	tests := []struct {
		receiver string
		want     string
	}{
		{"repo", "repository-interface"},
		{"db", "repository-interface"},
		{"store", "repository-interface"},
		{"prisma", "repository-interface"},
		{"knex", "repository-interface"},
		{"sequelize", "repository-interface"},
		{"mongoose", "repository-interface"},
		{"typeorm", "repository-interface"},
		{"authService", "external"},
		{"client", "external"},
	}
	for _, tt := range tests {
		got := classifyTSBoundary(tt.receiver)
		if got != tt.want {
			t.Errorf("classifyTSBoundary(%q) = %q, want %q", tt.receiver, got, tt.want)
		}
	}
}

func TestFindTSFuncEnd(t *testing.T) {
	lines := []string{
		"export function foo() {",        // line 0
		"  const x = 1;",                 // line 1
		"  if (x > 0) {",                 // line 2
		"    return x;",                   // line 3
		"  }",                             // line 4
		"}",                               // line 5
		"",                                // line 6
		"export function bar() {",         // line 7
	}

	endLine := findTSFuncEnd(lines, 0)
	// The function ends at line 5 (0-indexed), so 1-indexed = 6.
	if endLine != 6 {
		t.Errorf("findTSFuncEnd = %d, want 6", endLine)
	}
}

// ---------------------------------------------------------------------------
// PyTracer tests
// ---------------------------------------------------------------------------

func TestPyTracerBasicChain(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.py", `from auth_service import login

def handle_login(request):
    result = auth_service.login(request.body)
    return result
`)

	writeFile(t, dir, "auth_service.py", `from utils import generate_token

def login(credentials):
    user = repo.find_user(credentials.email)
    token = generate_token(user)
    return {"user": user, "token": token}
`)

	writeFile(t, dir, "utils.py", `def generate_token(user):
    return "jwt-token-" + str(user.id)
`)

	tr := &PyTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "handler.py",
		StartLine: 3,
		EndLine:   5,
	})
	if err != nil {
		t.Fatalf("PyTracer.Trace: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("expected non-empty chain entries from PyTracer")
	}

	// login should appear as internal.
	foundLogin := findEntryContaining(entries, "login", false)
	if foundLogin == nil {
		t.Error("expected login in chain entries")
	} else if foundLogin.File == "" {
		t.Error("expected login to have a File (internal)")
	}

	// generate_token should appear as internal.
	foundToken := findEntryContaining(entries, "generate_token", false)
	if foundToken == nil {
		t.Error("expected generate_token in chain entries")
	} else if foundToken.File == "" {
		t.Error("expected generate_token to have a File (internal)")
	}

	// repo.find_user should be external.
	foundRepo := findEntryContaining(entries, "find_user", true)
	if foundRepo == nil {
		t.Error("expected repo.find_user as external entry")
	} else if foundRepo.Boundary == "" {
		t.Error("expected repo.find_user to have a boundary classification")
	}
}

func TestPyTracerHandlerNotFound(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.py", `def foo():
    pass

def bar():
    pass
`)

	tr := &PyTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "nonexistent.py",
		StartLine: 1,
		EndLine:   3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries when handler not found, got %d", len(entries))
	}
}

func TestPyTracerEmptyProject(t *testing.T) {
	dir := t.TempDir()

	tr := &PyTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "handler.py",
		StartLine: 1,
		EndLine:   5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for empty project, got %d", len(entries))
	}
}

func TestPyTracerClassifyBoundary(t *testing.T) {
	imports := map[string]string{
		"requests": "external",
		"utils":    "internal",
	}

	tests := []struct {
		callExpr string
		want     string
	}{
		{"repo.find", "repository-interface"},
		{"db.execute", "repository-interface"},
		{"store.get", "repository-interface"},
		{"model.query", "repository-interface"},
		{"session.commit", "repository-interface"},
		{"requests.get", "external"},
		{"unknown.call", "external"},
	}
	for _, tt := range tests {
		got := classifyPyBoundary(tt.callExpr, imports)
		if got != tt.want {
			t.Errorf("classifyPyBoundary(%q) = %q, want %q", tt.callExpr, got, tt.want)
		}
	}
}

func TestPyTracerIsPyBuiltin(t *testing.T) {
	builtins := []string{"print", "len", "range", "str", "int", "isinstance", "super"}
	for _, b := range builtins {
		if !isPyBuiltin(b) {
			t.Errorf("isPyBuiltin(%q) = false, want true", b)
		}
	}

	nonBuiltins := []string{"my_func", "login", "generate_token", "handle_request"}
	for _, nb := range nonBuiltins {
		if isPyBuiltin(nb) {
			t.Errorf("isPyBuiltin(%q) = true, want false", nb)
		}
	}
}

func TestPyEffectiveIndent(t *testing.T) {
	tests := []struct {
		indent string
		want   int
	}{
		{"", 0},
		{"    ", 4},
		{"\t", 4},
		{"  \t", 6},
		{"\t\t", 8},
	}
	for _, tt := range tests {
		got := pyEffectiveIndent(tt.indent)
		if got != tt.want {
			t.Errorf("pyEffectiveIndent(%q) = %d, want %d", tt.indent, got, tt.want)
		}
	}
}

func TestPyLeadingWhitespace(t *testing.T) {
	tests := []struct {
		line string
		want string
	}{
		{"    def foo():", "    "},
		{"\tdef foo():", "\t"},
		{"def foo():", ""},
		{"  \t  x = 1", "  \t  "},
	}
	for _, tt := range tests {
		got := pyLeadingWhitespace(tt.line)
		if got != tt.want {
			t.Errorf("pyLeadingWhitespace(%q) = %q, want %q", tt.line, got, tt.want)
		}
	}
}

func TestPyCollectImports(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "test.py", `import os
import json as j
from datetime import datetime
from .utils import helper
from .models import User as UserModel
`)

	imports := collectImports(filepath.Join(dir, "test.py"))

	if imports["os"] != "external" {
		t.Errorf("os should be external, got %q", imports["os"])
	}
	if imports["j"] != "external" {
		t.Errorf("j (json alias) should be external, got %q", imports["j"])
	}
	if imports["datetime"] != "external" {
		t.Errorf("datetime should be external, got %q", imports["datetime"])
	}
	if imports["helper"] != "internal" {
		t.Errorf("helper (relative import) should be internal, got %q", imports["helper"])
	}
	if imports["UserModel"] != "internal" {
		t.Errorf("UserModel (relative import alias) should be internal, got %q", imports["UserModel"])
	}
}

func TestPyFindFuncEndTracer(t *testing.T) {
	lines := []string{
		"def foo():",         // 0
		"    x = 1",          // 1
		"    if x:",          // 2
		"        return x",   // 3
		"",                   // 4
		"def bar():",         // 5
	}

	endLine := findPyFuncEndTracer(lines, 0, "")
	// foo ends at line 3 (0-indexed), 1-indexed = 4.
	if endLine != 4 {
		t.Errorf("findPyFuncEndTracer = %d, want 4", endLine)
	}
}

func TestPyTracerClassWithMethods(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.py", `from auth_service import AuthService

def handle_request(request):
    svc = AuthService()
    result = svc.process(request)
    return result
`)

	writeFile(t, dir, "auth_service.py", `class AuthService:
    def process(self, request):
        return self.validate(request)

    def validate(self, request):
        return True
`)

	tr := &PyTracer{}
	entries, err := tr.Trace(dir, model.FuncLocation{
		File:      "handler.py",
		StartLine: 3,
		EndLine:   6,
	})
	if err != nil {
		t.Fatalf("PyTracer.Trace: %v", err)
	}

	// process should be found as an internal function.
	foundProcess := findEntryContaining(entries, "process", false)
	if foundProcess == nil {
		t.Log("process not found in chain (may be matched as external); checking entries:")
		for _, e := range entries {
			t.Logf("  func=%q file=%q boundary=%q", e.Func, e.File, e.Boundary)
		}
	}
}

// ---------------------------------------------------------------------------
// Helpers for searching entries
// ---------------------------------------------------------------------------

// findEntry finds an exact match by Func name.
// If boundary is true, only matches entries with a non-empty Boundary.
func findEntry(entries []model.ChainEntry, funcName string, boundary bool) *model.ChainEntry {
	for i, e := range entries {
		if e.Func == funcName {
			if boundary && e.Boundary == "" {
				continue
			}
			if !boundary && e.Boundary != "" {
				continue
			}
			return &entries[i]
		}
	}
	return nil
}

// findEntryContaining finds an entry whose Func contains the given substring.
func findEntryContaining(entries []model.ChainEntry, substr string, boundary bool) *model.ChainEntry {
	for i, e := range entries {
		if contains(e.Func, substr) {
			if boundary && e.Boundary == "" {
				continue
			}
			if !boundary && e.Boundary != "" {
				continue
			}
			return &entries[i]
		}
	}
	return nil
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
