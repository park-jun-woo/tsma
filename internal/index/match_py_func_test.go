package index

import "testing"

func TestMatchPyFuncTopLevel(t *testing.T) {
	fn, newClass := matchPyFunc("def handle_login(request):", 5, "handler.py", "", "", 0)
	if fn == nil {
		t.Fatal("expected match for top-level function")
	}
	if fn.Name != "handle_login" {
		t.Errorf("Name = %q, want %q", fn.Name, "handle_login")
	}
	if fn.QualifiedName != "handle_login" {
		t.Errorf("QualifiedName = %q, want %q", fn.QualifiedName, "handle_login")
	}
	if fn.StartLine != 5 {
		t.Errorf("StartLine = %d, want 5", fn.StartLine)
	}
	if fn.Exported {
		t.Error("expected Exported=false for lowercase function")
	}
	if newClass != "" {
		t.Errorf("newClass = %q, want empty", newClass)
	}
}

func TestMatchPyFuncMethod(t *testing.T) {
	fn, newClass := matchPyFunc("    def login(self, credentials):", 10, "services/auth.py", "services", "AuthService", 0)
	if fn == nil {
		t.Fatal("expected match for class method")
	}
	if fn.Name != "login" {
		t.Errorf("Name = %q, want %q", fn.Name, "login")
	}
	if fn.QualifiedName != "services.AuthService.login" {
		t.Errorf("QualifiedName = %q, want %q", fn.QualifiedName, "services.AuthService.login")
	}
	if newClass != "AuthService" {
		t.Errorf("newClass = %q, want %q", newClass, "AuthService")
	}
}

func TestMatchPyFuncResetsClass(t *testing.T) {
	// A function at indent 0 when class was at indent 0 should reset the class context
	fn, newClass := matchPyFunc("def standalone():", 20, "handler.py", "", "MyClass", 0)
	if fn == nil {
		t.Fatal("expected match for standalone function")
	}
	if fn.Name != "standalone" {
		t.Errorf("Name = %q, want %q", fn.Name, "standalone")
	}
	// Class should be reset since indent <= classIndent
	if newClass != "" {
		t.Errorf("newClass = %q, want empty (class context should be reset)", newClass)
	}
}

func TestMatchPyFuncAsyncDef(t *testing.T) {
	fn, _ := matchPyFunc("async def fetch_data(url):", 1, "client.py", "", "", 0)
	if fn == nil {
		t.Fatal("expected match for async def")
	}
	if fn.Name != "fetch_data" {
		t.Errorf("Name = %q, want %q", fn.Name, "fetch_data")
	}
}

func TestMatchPyFuncNoMatch(t *testing.T) {
	fn, newClass := matchPyFunc("x = 42", 1, "handler.py", "", "", 0)
	if fn != nil {
		t.Error("expected nil for non-function line")
	}
	if newClass != "" {
		t.Errorf("newClass = %q, want empty", newClass)
	}
}

func TestMatchPyFuncExported(t *testing.T) {
	fn, _ := matchPyFunc("def HandleLogin(request):", 1, "handler.py", "", "", 0)
	if fn == nil {
		t.Fatal("expected match")
	}
	if !fn.Exported {
		t.Error("expected Exported=true for uppercase function")
	}
}

func TestMatchPyFuncPrivate(t *testing.T) {
	fn, _ := matchPyFunc("def _private_helper():", 1, "handler.py", "", "", 0)
	if fn == nil {
		t.Fatal("expected match")
	}
	if fn.Exported {
		t.Error("expected Exported=false for _private function")
	}
}
