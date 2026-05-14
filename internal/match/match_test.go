package match

import (
	"os"
	"path/filepath"
	"testing"
)

// writeFixture creates a file relative to dir with the given content.
func writeFixture(t *testing.T, dir, rel, content string) {
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
// Factory tests
// ---------------------------------------------------------------------------

func TestNewMatcherGo(t *testing.T) {
	m := NewMatcher("go")
	if _, ok := m.(*GoMatcher); !ok {
		t.Errorf("NewMatcher(\"go\") returned %T, want *GoMatcher", m)
	}
}

func TestNewMatcherTypescript(t *testing.T) {
	m := NewMatcher("typescript")
	if _, ok := m.(*TSMatcher); !ok {
		t.Errorf("NewMatcher(\"typescript\") returned %T, want *TSMatcher", m)
	}
}

func TestNewMatcherPython(t *testing.T) {
	m := NewMatcher("python")
	if _, ok := m.(*PyMatcher); !ok {
		t.Errorf("NewMatcher(\"python\") returned %T, want *PyMatcher", m)
	}
}

func TestNewMatcherUnsupported(t *testing.T) {
	m := NewMatcher("ruby")
	_, found := m.Match("/tmp", "handler.rb")
	if found {
		t.Error("unsupported matcher should return found=false")
	}
}

// ---------------------------------------------------------------------------
// GoMatcher tests
// ---------------------------------------------------------------------------

func TestGoMatcherFound(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "internal/handler/handler.go", `package handler

func Login() error { return nil }
`)
	writeFixture(t, dir, "internal/handler/handler_test.go", `package handler

import "testing"

func TestLogin(t *testing.T) {}
`)

	m := &GoMatcher{}
	testFile, found := m.Match(dir, "internal/handler/handler.go")
	if !found {
		t.Fatal("expected to find test file for handler.go")
	}
	if testFile != filepath.Join("internal", "handler", "handler_test.go") {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestGoMatcherNotFound(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "handler.go", `package main

func Login() error { return nil }
`)

	m := &GoMatcher{}
	_, found := m.Match(dir, "handler.go")
	if found {
		t.Error("expected no match when no test file exists")
	}
}

func TestGoMatcherNoTestFiles(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "internal/api/server.go", `package api

func Serve() {}
`)

	m := &GoMatcher{}
	_, found := m.Match(dir, "internal/api/server.go")
	if found {
		t.Error("expected no match when no test file exists")
	}
}

// ---------------------------------------------------------------------------
// TSMatcher tests
// ---------------------------------------------------------------------------

func TestTSMatcherTestFile(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "src/handler.ts", `export function handleLogin() {}`)
	writeFixture(t, dir, "src/handler.test.ts", `describe('handleLogin', () => {});`)

	m := &TSMatcher{}
	testFile, found := m.Match(dir, "src/handler.ts")
	if !found {
		t.Fatal("expected to find test file for handler.ts")
	}
	if testFile != filepath.Join("src", "handler.test.ts") {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestTSMatcherSpecFile(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "src/service.ts", `export function createUser() {}`)
	writeFixture(t, dir, "src/service.spec.ts", `describe('createUser', () => {});`)

	m := &TSMatcher{}
	testFile, found := m.Match(dir, "src/service.ts")
	if !found {
		t.Fatal("expected to find .spec.ts file")
	}
	if testFile != filepath.Join("src", "service.spec.ts") {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestTSMatcherTestsDirMatch(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "src/handler.ts", `export function handleLogin() {}`)
	writeFixture(t, dir, "src/__tests__/handler.test.ts", `describe('handleLogin', () => {});`)

	m := &TSMatcher{}
	testFile, found := m.Match(dir, "src/handler.ts")
	if !found {
		t.Fatal("expected to find test file in __tests__/ dir")
	}
	if testFile != filepath.Join("src", "__tests__", "handler.test.ts") {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestTSMatcherNotFound(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "src/handler.ts", `export function handleLogin() {}`)

	m := &TSMatcher{}
	_, found := m.Match(dir, "src/handler.ts")
	if found {
		t.Error("expected no match when no test files exist")
	}
}

// ---------------------------------------------------------------------------
// PyMatcher tests
// ---------------------------------------------------------------------------

func TestPyMatcherFilenameMatch(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "auth_svc.py", `def login(): pass`)
	writeFixture(t, dir, "test_auth_svc.py", `def test_login(): pass`)

	m := &PyMatcher{}
	testFile, found := m.Match(dir, "auth_svc.py")
	if !found {
		t.Fatal("expected to find test file for auth_svc.py")
	}
	if testFile != "test_auth_svc.py" {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestPyMatcherTestsDirMatch(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "src/handler.py", `def handle_request(): pass`)
	writeFixture(t, dir, "src/tests/test_handler.py", `def test_handle_request(): pass`)

	m := &PyMatcher{}
	testFile, found := m.Match(dir, "src/handler.py")
	if !found {
		t.Fatal("expected to find test file in tests/ dir")
	}
	if testFile != filepath.Join("src", "tests", "test_handler.py") {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestPyMatcherNotFound(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "handler.py", `def login(): pass`)

	m := &PyMatcher{}
	_, found := m.Match(dir, "handler.py")
	if found {
		t.Error("expected no match when no test files exist")
	}
}

// ---------------------------------------------------------------------------
// Helper function tests
// ---------------------------------------------------------------------------

func TestStripTSExtension(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"handler.ts", "handler"},
		{"handler.js", "handler"},
		{"handler.tsx", "handler"},
		{"handler.jsx", "handler"},
		{"handler.py", "handler.py"},
		{"handler", "handler"},
	}
	for _, tt := range tests {
		got := stripTSExtension(tt.input)
		if got != tt.want {
			t.Errorf("stripTSExtension(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsTSTestFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"handler.test.ts", true},
		{"handler.test.js", true},
		{"handler.spec.ts", true},
		{"handler.spec.js", true},
		{"handler.ts", false},
		{"handler.go", false},
	}
	for _, tt := range tests {
		got := isTSTestFile(tt.name)
		if got != tt.want {
			t.Errorf("isTSTestFile(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}
