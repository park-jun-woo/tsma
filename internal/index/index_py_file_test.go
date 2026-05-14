package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIndexPyFileTopLevel(t *testing.T) {
	dir := t.TempDir()
	content := `def handle_login(request):
    result = auth_service.login(request.body)
    return result

async def async_handler(request):
    return await process(request)
`
	absPath := filepath.Join(dir, "handler.py")
	if err := os.WriteFile(absPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs := indexPyFile("handler.py", absPath)
	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d", len(funcs))
	}

	names := map[string]bool{}
	for _, fn := range funcs {
		names[fn.Name] = true
	}
	if !names["handle_login"] {
		t.Error("expected handle_login")
	}
	if !names["async_handler"] {
		t.Error("expected async_handler")
	}
}

func TestIndexPyFileWithClass(t *testing.T) {
	dir := t.TempDir()
	content := `class AuthService:
    def login(self, credentials):
        return self.validate(credentials)

    def validate(self, data):
        return True

def standalone():
    pass
`
	absPath := filepath.Join(dir, "services", "auth.py")
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs := indexPyFile("services/auth.py", absPath)
	if len(funcs) < 3 {
		t.Fatalf("expected at least 3 functions, got %d", len(funcs))
	}

	var loginFound bool
	for _, fn := range funcs {
		if fn.Name == "login" {
			loginFound = true
			if fn.QualifiedName != "services.AuthService.login" {
				t.Errorf("login QualifiedName = %q, want %q", fn.QualifiedName, "services.AuthService.login")
			}
		}
	}
	if !loginFound {
		t.Error("expected to find login method")
	}
}

func TestIndexPyFileNonExistent(t *testing.T) {
	funcs := indexPyFile("nonexistent.py", "/nonexistent/path/handler.py")
	if funcs != nil {
		t.Errorf("expected nil for non-existent file, got %v", funcs)
	}
}
