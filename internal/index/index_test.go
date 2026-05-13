package index

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/model"
)

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
// Factory tests
// ---------------------------------------------------------------------------

func TestNewIndexerGo(t *testing.T) {
	idx := NewIndexer("go")
	if _, ok := idx.(*GoIndexer); !ok {
		t.Errorf("NewIndexer(\"go\") returned %T, want *GoIndexer", idx)
	}
}

func TestNewIndexerTypescript(t *testing.T) {
	idx := NewIndexer("typescript")
	if _, ok := idx.(*TSIndexer); !ok {
		t.Errorf("NewIndexer(\"typescript\") returned %T, want *TSIndexer", idx)
	}
}

func TestNewIndexerPython(t *testing.T) {
	idx := NewIndexer("python")
	if _, ok := idx.(*PyIndexer); !ok {
		t.Errorf("NewIndexer(\"python\") returned %T, want *PyIndexer", idx)
	}
}

func TestNewIndexerUnsupported(t *testing.T) {
	idx := NewIndexer("ruby")
	u, ok := idx.(*UnsupportedIndexer)
	if !ok {
		t.Fatalf("NewIndexer(\"ruby\") returned %T, want *UnsupportedIndexer", idx)
	}
	if u.Lang != "ruby" {
		t.Errorf("UnsupportedIndexer.Lang = %q, want \"ruby\"", u.Lang)
	}
	_, err := idx.Index(".")
	if err == nil {
		t.Fatal("expected error from UnsupportedIndexer.Index")
	}
}

// ---------------------------------------------------------------------------
// GoIndexer tests
// ---------------------------------------------------------------------------

func TestGoIndexerBasic(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.go", `package main

func main() {
	handleLogin()
}
`)

	writeFile(t, dir, "handler.go", `package main

func handleLogin() {
	svc.Login()
}
`)

	idx := &GoIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("GoIndexer.Index: %v", err)
	}

	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d", len(funcs))
	}

	found := findByName(funcs, "main")
	if found == nil {
		t.Error("expected to find function 'main'")
	}

	found = findByName(funcs, "handleLogin")
	if found == nil {
		t.Error("expected to find function 'handleLogin'")
	} else {
		if found.Exported {
			t.Error("handleLogin should not be exported")
		}
		if found.IsMethod {
			t.Error("handleLogin should not be a method")
		}
	}
}

func TestGoIndexerWithReceiver(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "internal/api/handler.go", `package api

type Handler struct{}

func (h *Handler) Login() error {
	return nil
}

func (h Handler) Logout() error {
	return nil
}

func helperFunc() string {
	return "ok"
}
`)

	idx := &GoIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("GoIndexer.Index: %v", err)
	}

	if len(funcs) != 3 {
		t.Fatalf("expected 3 functions, got %d", len(funcs))
	}

	login := findByQualifiedName(funcs, "internal/api.Handler.Login")
	if login == nil {
		t.Fatal("expected to find 'internal/api.Handler.Login'")
	}
	if !login.IsMethod {
		t.Error("Login should be a method")
	}
	if login.Receiver != "Handler" {
		t.Errorf("Login receiver = %q, want \"Handler\"", login.Receiver)
	}
	if !login.Exported {
		t.Error("Login should be exported")
	}

	logout := findByQualifiedName(funcs, "internal/api.Handler.Logout")
	if logout == nil {
		t.Fatal("expected to find 'internal/api.Handler.Logout'")
	}

	helper := findByQualifiedName(funcs, "internal/api.helperFunc")
	if helper == nil {
		t.Fatal("expected to find 'internal/api.helperFunc'")
	}
	if helper.IsMethod {
		t.Error("helperFunc should not be a method")
	}
	if helper.Exported {
		t.Error("helperFunc should not be exported")
	}
}

func TestGoIndexerSkipsTestFiles(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.go", `package main

func Handler() {}
`)

	writeFile(t, dir, "handler_test.go", `package main

import "testing"

func TestHandler(t *testing.T) {}
`)

	writeFile(t, dir, "mock_service.go", `package main

func MockLogin() {}
`)

	idx := &GoIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("GoIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (test and mock files excluded), got %d", len(funcs))
	}

	if funcs[0].Name != "Handler" {
		t.Errorf("expected function 'Handler', got %q", funcs[0].Name)
	}
}

func TestGoIndexerSkipsVendor(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.go", `package main

func main() {}
`)
	writeFile(t, dir, "vendor/lib/lib.go", `package lib

func Vendored() {}
`)

	idx := &GoIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("GoIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (vendor excluded), got %d", len(funcs))
	}
}

func TestGoIndexerEmptyProject(t *testing.T) {
	dir := t.TempDir()

	idx := &GoIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("GoIndexer.Index: %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("expected 0 functions for empty project, got %d", len(funcs))
	}
}

func TestGoIndexerLineNumbers(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.go", `package main

func first() {
	// line 4
}

func second() {
	// line 8
}
`)

	idx := &GoIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("GoIndexer.Index: %v", err)
	}

	first := findByName(funcs, "first")
	if first == nil {
		t.Fatal("expected to find 'first'")
	}
	if first.StartLine != 3 {
		t.Errorf("first.StartLine = %d, want 3", first.StartLine)
	}

	second := findByName(funcs, "second")
	if second == nil {
		t.Fatal("expected to find 'second'")
	}
	if second.StartLine != 7 {
		t.Errorf("second.StartLine = %d, want 7", second.StartLine)
	}
}

// ---------------------------------------------------------------------------
// TSIndexer tests
// ---------------------------------------------------------------------------

func TestTSIndexerBasic(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.ts", `export async function handleLogin(req: Request) {
  const result = await authService.login(req.body);
  return result;
}

function helperFunc() {
  return "ok";
}

export const formatDate = (d: Date) => {
  return d.toISOString();
};
`)

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}

	if len(funcs) < 3 {
		t.Fatalf("expected at least 3 functions, got %d", len(funcs))
	}

	hl := findByName(funcs, "handleLogin")
	if hl == nil {
		t.Error("expected to find 'handleLogin'")
	} else if !hl.Exported {
		t.Error("handleLogin should be exported")
	}

	helper := findByName(funcs, "helperFunc")
	if helper == nil {
		t.Error("expected to find 'helperFunc'")
	} else if helper.Exported {
		t.Error("helperFunc should not be exported")
	}

	format := findByName(funcs, "formatDate")
	if format == nil {
		t.Error("expected to find 'formatDate'")
	} else if !format.Exported {
		t.Error("formatDate should be exported")
	}
}

func TestTSIndexerClassMethods(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "src/auth/service.ts", `export class AuthService {
  login(credentials: any) {
    return this.validate(credentials);
  }

  async validate(data: any) {
    return true;
  }
}
`)

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}

	login := findByName(funcs, "login")
	if login == nil {
		t.Fatal("expected to find 'login'")
	}
	if !login.IsMethod {
		t.Error("login should be a method")
	}
	if login.Receiver != "AuthService" {
		t.Errorf("login.Receiver = %q, want \"AuthService\"", login.Receiver)
	}
	if login.QualifiedName != "src/auth.AuthService.login" {
		t.Errorf("login.QualifiedName = %q, want \"src/auth.AuthService.login\"", login.QualifiedName)
	}

	validate := findByName(funcs, "validate")
	if validate == nil {
		t.Fatal("expected to find 'validate'")
	}
}

func TestTSIndexerSkipsTestFiles(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.ts", `export function handler() {}
`)
	writeFile(t, dir, "handler.test.ts", `describe('handler', () => {});
`)
	writeFile(t, dir, "handler.spec.ts", `describe('handler', () => {});
`)
	writeFile(t, dir, "types.d.ts", `declare function foo(): void;
`)

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (test/spec/d.ts excluded), got %d", len(funcs))
	}
}

func TestTSIndexerSkipsNodeModules(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "index.ts", `export function main() {}
`)
	writeFile(t, dir, "node_modules/lib/index.ts", `export function libFunc() {}
`)

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (node_modules excluded), got %d", len(funcs))
	}
}

func TestTSIndexerSubdirectory(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "src/api/auth.ts", `export function login() {}
`)

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function, got %d", len(funcs))
	}

	if funcs[0].QualifiedName != "src/api.login" {
		t.Errorf("qualified_name = %q, want \"src/api.login\"", funcs[0].QualifiedName)
	}
	if funcs[0].File != "src/api/auth.ts" {
		t.Errorf("file = %q, want \"src/api/auth.ts\"", funcs[0].File)
	}
}

// ---------------------------------------------------------------------------
// PyIndexer tests
// ---------------------------------------------------------------------------

func TestPyIndexerBasic(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.py", `def handle_login(request):
    result = auth_service.login(request.body)
    return result

async def async_handler(request):
    return await process(request)
`)

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}

	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d", len(funcs))
	}

	hl := findByName(funcs, "handle_login")
	if hl == nil {
		t.Error("expected to find 'handle_login'")
	}

	ah := findByName(funcs, "async_handler")
	if ah == nil {
		t.Error("expected to find 'async_handler'")
	}
}

func TestPyIndexerClassMethods(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "services/auth.py", `class AuthService:
    def login(self, credentials):
        return self.validate(credentials)

    def validate(self, data):
        return True

def standalone():
    pass
`)

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}

	if len(funcs) < 3 {
		t.Fatalf("expected at least 3 functions, got %d", len(funcs))
	}

	login := findByName(funcs, "login")
	if login == nil {
		t.Fatal("expected to find 'login'")
	}
	if !login.IsMethod {
		t.Error("login should be a method")
	}
	if login.Receiver != "AuthService" {
		t.Errorf("login.Receiver = %q, want \"AuthService\"", login.Receiver)
	}
	if login.QualifiedName != "services.AuthService.login" {
		t.Errorf("login.QualifiedName = %q, want \"services.AuthService.login\"", login.QualifiedName)
	}

	standalone := findByName(funcs, "standalone")
	if standalone == nil {
		t.Fatal("expected to find 'standalone'")
	}
	if standalone.IsMethod {
		t.Error("standalone should not be a method")
	}
}

func TestPyIndexerSkipsTestFiles(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "handler.py", `def handler():
    pass
`)
	writeFile(t, dir, "test_handler.py", `def test_handler():
    pass
`)

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (test files excluded), got %d", len(funcs))
	}
}

func TestPyIndexerSkipsPycache(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "main.py", `def main():
    pass
`)
	writeFile(t, dir, "__pycache__/cached.py", `def cached():
    pass
`)
	writeFile(t, dir, ".venv/lib/lib.py", `def libfunc():
    pass
`)

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (__pycache__/.venv excluded), got %d", len(funcs))
	}
}

func TestPyIndexerSubdirectory(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "app/api/auth.py", `def login():
    pass
`)

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function, got %d", len(funcs))
	}

	if funcs[0].QualifiedName != "app/api.login" {
		t.Errorf("qualified_name = %q, want \"app/api.login\"", funcs[0].QualifiedName)
	}
}

// ---------------------------------------------------------------------------
// Helper filter tests
// ---------------------------------------------------------------------------

func TestIsGoSource(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"handler.go", true},
		{"handler_test.go", false},
		{"mock_service.go", false},
		{"readme.md", false},
		{"internal/api/handler.go", true},
		{"internal/api/handler_test.go", false},
		{"internal/mock_repo.go", false},
	}
	for _, tt := range tests {
		got := isGoSource(tt.path)
		if got != tt.want {
			t.Errorf("isGoSource(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestIsTSSource(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"handler.ts", true},
		{"service.js", true},
		{"types.d.ts", false},
		{"handler.test.ts", false},
		{"handler.spec.ts", false},
		{"handler.test.js", false},
		{"handler.spec.js", false},
		{"readme.md", false},
	}
	for _, tt := range tests {
		got := isTSSource(tt.path)
		if got != tt.want {
			t.Errorf("isTSSource(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestIsPySource(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"handler.py", true},
		{"test_handler.py", false},
		{"readme.md", false},
		{"app/auth.py", true},
		{"app/test_auth.py", false},
	}
	for _, tt := range tests {
		got := isPySource(tt.path)
		if got != tt.want {
			t.Errorf("isPySource(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Helpers for searching functions
// ---------------------------------------------------------------------------

func findByName(funcs []model.Function, name string) *model.Function {
	for i := range funcs {
		if funcs[i].Name == name {
			return &funcs[i]
		}
	}
	return nil
}

func findByQualifiedName(funcs []model.Function, qn string) *model.Function {
	for i := range funcs {
		if funcs[i].QualifiedName == qn {
			return &funcs[i]
		}
	}
	return nil
}
