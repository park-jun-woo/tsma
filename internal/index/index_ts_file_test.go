package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIndexTSFileTopLevel(t *testing.T) {
	dir := t.TempDir()
	content := `export async function handleLogin(req: Request) {
  const result = await authService.login(req.body);
  return result;
}

function helperFunc() {
  return "ok";
}

export const formatDate = (d: Date) => {
  return d.toISOString();
};
`
	absPath := filepath.Join(dir, "handler.ts")
	if err := os.WriteFile(absPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs := indexTSFile("handler.ts", absPath)
	if len(funcs) < 3 {
		t.Fatalf("expected at least 3 functions, got %d", len(funcs))
	}

	names := map[string]bool{}
	for _, fn := range funcs {
		names[fn.Name] = true
	}
	if !names["handleLogin"] {
		t.Error("expected handleLogin")
	}
	if !names["helperFunc"] {
		t.Error("expected helperFunc")
	}
	if !names["formatDate"] {
		t.Error("expected formatDate")
	}
}

func TestIndexTSFileWithClass(t *testing.T) {
	dir := t.TempDir()
	content := `export class AuthService {
  login(credentials: any) {
    return this.validate(credentials);
  }

  async validate(data: any) {
    return true;
  }
}
`
	absPath := filepath.Join(dir, "src", "auth", "service.ts")
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs := indexTSFile("src/auth/service.ts", absPath)
	if len(funcs) < 2 {
		t.Fatalf("expected at least 2 functions, got %d", len(funcs))
	}

	var loginFound bool
	for _, fn := range funcs {
		if fn.Name == "login" {
			loginFound = true
			if fn.QualifiedName != "src/auth.AuthService.login" {
				t.Errorf("login QualifiedName = %q, want %q", fn.QualifiedName, "src/auth.AuthService.login")
			}
		}
	}
	if !loginFound {
		t.Error("expected to find login method")
	}
}

func TestIndexTSFileClassReset(t *testing.T) {
	dir := t.TempDir()
	content := `class MyClass {
  doSomething() {
    return 1;
  }
}

function standalone() {
  return 2;
}
`
	absPath := filepath.Join(dir, "mixed.ts")
	if err := os.WriteFile(absPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs := indexTSFile("mixed.ts", absPath)

	var standaloneFound bool
	for _, fn := range funcs {
		if fn.Name == "standalone" {
			standaloneFound = true
			// After class closes, standalone should not have a class receiver
			if fn.QualifiedName != "standalone" {
				t.Errorf("standalone QualifiedName = %q, want %q", fn.QualifiedName, "standalone")
			}
		}
	}
	if !standaloneFound {
		t.Error("expected to find standalone function")
	}
}

func TestIndexTSFileNonExistent(t *testing.T) {
	funcs := indexTSFile("nonexistent.ts", "/nonexistent/path/handler.ts")
	if funcs != nil {
		t.Errorf("expected nil for non-existent file, got %v", funcs)
	}
}
