package index

import "testing"

func TestMatchTSMethodBasic(t *testing.T) {
	fn, ok := matchTSMethod("  login(credentials: any) {", "AuthService", "src/auth", "src/auth/service.ts", 10)
	if !ok {
		t.Fatal("expected match for class method")
	}
	if fn.Name != "login" {
		t.Errorf("Name = %q, want %q", fn.Name, "login")
	}
	if fn.QualifiedName != "src/auth.AuthService.login" {
		t.Errorf("QualifiedName = %q, want %q", fn.QualifiedName, "src/auth.AuthService.login")
	}
	if fn.StartLine != 10 {
		t.Errorf("StartLine = %d, want 10", fn.StartLine)
	}
}

func TestMatchTSMethodAsync(t *testing.T) {
	fn, ok := matchTSMethod("  async validate(data: any): Promise<boolean> {", "AuthService", "src/auth", "src/auth/service.ts", 15)
	if !ok {
		t.Fatal("expected match for async method")
	}
	if fn.Name != "validate" {
		t.Errorf("Name = %q, want %q", fn.Name, "validate")
	}
}

func TestMatchTSMethodSkipsConstructor(t *testing.T) {
	_, ok := matchTSMethod("  constructor(private svc: Service) {", "MyClass", "", "file.ts", 1)
	if ok {
		t.Error("expected constructor to be skipped")
	}
}

func TestMatchTSMethodSkipsKeywords(t *testing.T) {
	keywords := []string{
		"  if (condition) {",
		"  for (let i = 0) {",
		"  while (true) {",
		"  switch (val) {",
	}
	for _, line := range keywords {
		_, ok := matchTSMethod(line, "MyClass", "", "file.ts", 1)
		if ok {
			t.Errorf("expected skip for keyword line: %q", line)
		}
	}
}

func TestMatchTSMethodNoMatch(t *testing.T) {
	_, ok := matchTSMethod("const x = 42;", "MyClass", "", "file.ts", 1)
	if ok {
		t.Error("expected no match for non-method line")
	}
}
