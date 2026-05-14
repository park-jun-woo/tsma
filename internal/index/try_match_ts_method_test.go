package index

import "testing"

func TestTryMatchTSMethodWithClass(t *testing.T) {
	fn, ok := tryMatchTSMethod("  login(credentials: any) {", "AuthService", "src/auth", "src/auth/service.ts", 10)
	if !ok {
		t.Fatal("expected match when currentClass is set")
	}
	if fn.Name != "login" {
		t.Errorf("Name = %q, want %q", fn.Name, "login")
	}
}

func TestTryMatchTSMethodWithoutClass(t *testing.T) {
	_, ok := tryMatchTSMethod("  login(credentials: any) {", "", "src/auth", "src/auth/service.ts", 10)
	if ok {
		t.Error("expected no match when currentClass is empty")
	}
}
