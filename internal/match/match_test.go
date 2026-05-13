package match

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
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
	_, found := m.Match("/tmp", &model.Function{Name: "foo"})
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

func TestLogin_Success(t *testing.T) {}
func TestLogin_Failure(t *testing.T) {}
`)

	m := &GoMatcher{}
	fn := &model.Function{
		Name: "Login",
		File: "internal/handler/handler.go",
	}

	testFile, found := m.Match(dir, fn)
	if !found {
		t.Fatal("expected to find test file for Login")
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

	writeFixture(t, dir, "handler_test.go", `package main

import "testing"

func TestSignup(t *testing.T) {}
`)

	m := &GoMatcher{}
	fn := &model.Function{
		Name: "Login",
		File: "handler.go",
	}

	_, found := m.Match(dir, fn)
	if found {
		t.Error("expected no match when no Test* contains function name")
	}
}

func TestGoMatcherNoTestFiles(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "handler.go", `package main

func Login() error { return nil }
`)

	m := &GoMatcher{}
	fn := &model.Function{
		Name: "Login",
		File: "handler.go",
	}

	_, found := m.Match(dir, fn)
	if found {
		t.Error("expected no match when no test files exist")
	}
}

// ---------------------------------------------------------------------------
// TSMatcher tests
// ---------------------------------------------------------------------------

func TestTSMatcherFilenameMatch(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "src/handler.ts", `export function handleLogin() {}`)
	writeFixture(t, dir, "src/handler.test.ts", `describe('handleLogin', () => {
  test('should work', () => {});
});`)

	m := &TSMatcher{}
	fn := &model.Function{
		Name: "handleLogin",
		File: "src/handler.ts",
	}

	testFile, found := m.Match(dir, fn)
	if !found {
		t.Fatal("expected to find test file via filename match")
	}
	if testFile != filepath.Join("src", "handler.test.ts") {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestTSMatcherTestsDirMatch(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "src/handler.ts", `export function handleLogin() {}`)
	writeFixture(t, dir, "src/__tests__/handler.test.ts", `describe('handleLogin', () => {});`)

	m := &TSMatcher{}
	fn := &model.Function{
		Name: "handleLogin",
		File: "src/handler.ts",
	}

	testFile, found := m.Match(dir, fn)
	if !found {
		t.Fatal("expected to find test file in __tests__/ dir")
	}
	if testFile != filepath.Join("src", "__tests__", "handler.test.ts") {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestTSMatcherContentMatch(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "src/auth.ts", `export function login() {}`)
	// Different filename but content references the function.
	writeFixture(t, dir, "src/api.test.ts", `describe('login', () => {
  it('should authenticate', () => {});
});`)

	m := &TSMatcher{}
	fn := &model.Function{
		Name: "login",
		File: "src/auth.ts",
	}

	testFile, found := m.Match(dir, fn)
	if !found {
		t.Fatal("expected to find test file via content match")
	}
	if testFile != filepath.Join("src", "api.test.ts") {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestTSMatcherNotFound(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "src/handler.ts", `export function handleLogin() {}`)

	m := &TSMatcher{}
	fn := &model.Function{
		Name: "handleLogin",
		File: "src/handler.ts",
	}

	_, found := m.Match(dir, fn)
	if found {
		t.Error("expected no match when no test files exist")
	}
}

func TestTSMatcherSpecFile(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "src/service.ts", `export function createUser() {}`)
	writeFixture(t, dir, "src/service.spec.ts", `describe('createUser', () => {});`)

	m := &TSMatcher{}
	fn := &model.Function{
		Name: "createUser",
		File: "src/service.ts",
	}

	testFile, found := m.Match(dir, fn)
	if !found {
		t.Fatal("expected to find .spec.ts file")
	}
	if testFile != filepath.Join("src", "service.spec.ts") {
		t.Errorf("testFile = %q, unexpected", testFile)
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
	fn := &model.Function{
		Name: "login",
		File: "auth_svc.py",
	}

	testFile, found := m.Match(dir, fn)
	if !found {
		t.Fatal("expected to find test file via filename match")
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
	fn := &model.Function{
		Name: "handle_request",
		File: "src/handler.py",
	}

	testFile, found := m.Match(dir, fn)
	if !found {
		t.Fatal("expected to find test file in tests/ dir")
	}
	if testFile != filepath.Join("src", "tests", "test_handler.py") {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestPyMatcherContentMatch(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "service.py", `def create_order(): pass`)
	// Different filename but content mentions the function.
	writeFixture(t, dir, "test_orders.py", `def test_create_order():
    result = create_order()
    assert result is not None
`)

	m := &PyMatcher{}
	fn := &model.Function{
		Name: "create_order",
		File: "service.py",
	}

	testFile, found := m.Match(dir, fn)
	if !found {
		t.Fatal("expected to find test file via content match")
	}
	if testFile != "test_orders.py" {
		t.Errorf("testFile = %q, unexpected", testFile)
	}
}

func TestPyMatcherNotFound(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "handler.py", `def login(): pass`)

	m := &PyMatcher{}
	fn := &model.Function{
		Name: "login",
		File: "handler.py",
	}

	_, found := m.Match(dir, fn)
	if found {
		t.Error("expected no match when no test files exist")
	}
}

func TestPyMatcherNoContentMatch(t *testing.T) {
	dir := t.TempDir()

	writeFixture(t, dir, "handler.py", `def login(): pass`)
	writeFixture(t, dir, "test_other.py", `def test_signup(): pass`)

	m := &PyMatcher{}
	fn := &model.Function{
		Name: "login",
		File: "handler.py",
	}

	_, found := m.Match(dir, fn)
	if found {
		t.Error("expected no match when test file does not mention function")
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
