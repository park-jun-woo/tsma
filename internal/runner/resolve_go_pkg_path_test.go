package runner

import (
	"path/filepath"
	"testing"
)

func TestResolveGoPkgPath(t *testing.T) {
	projectRoot := "/home/user/myproject"
	absTest := "/home/user/myproject/internal/handler/handler_test.go"

	pkg, err := resolveGoPkgPath(projectRoot, absTest)
	if err != nil {
		t.Fatalf("resolveGoPkgPath: %v", err)
	}
	want := "./internal/handler"
	if pkg != want {
		t.Errorf("pkg = %q, want %q", pkg, want)
	}
}

func TestResolveGoPkgPathRootPackage(t *testing.T) {
	projectRoot := "/home/user/myproject"
	absTest := "/home/user/myproject/main_test.go"

	pkg, err := resolveGoPkgPath(projectRoot, absTest)
	if err != nil {
		t.Fatalf("resolveGoPkgPath: %v", err)
	}
	want := "./."
	if pkg != want {
		t.Errorf("pkg = %q, want %q", pkg, want)
	}
}

func TestResolveGoPkgPathNestedPackage(t *testing.T) {
	projectRoot := "/home/user/myproject"
	absTest := "/home/user/myproject/pkg/auth/middleware/middleware_test.go"

	pkg, err := resolveGoPkgPath(projectRoot, absTest)
	if err != nil {
		t.Fatalf("resolveGoPkgPath: %v", err)
	}
	want := "./pkg/auth/middleware"
	if pkg != want {
		t.Errorf("pkg = %q, want %q", pkg, want)
	}
}

func TestResolveGoPkgPathRelError(t *testing.T) {
	// An absolute project root and a relative test path cannot be made relative
	// to each other -> filepath.Rel returns an error.
	_, err := resolveGoPkgPath("/home/user/myproject", "relative/handler_test.go")
	if err == nil {
		t.Fatal("expected error when paths cannot be made relative")
	}
}

func TestResolveGoPkgPathUsesForwardSlash(t *testing.T) {
	// Even on systems that use backslash, result should use forward slash.
	projectRoot := "/home/user/myproject"
	absTest := filepath.Join(projectRoot, "internal", "handler", "handler_test.go")

	pkg, err := resolveGoPkgPath(projectRoot, absTest)
	if err != nil {
		t.Fatalf("resolveGoPkgPath: %v", err)
	}
	want := "./internal/handler"
	if pkg != want {
		t.Errorf("pkg = %q, want %q", pkg, want)
	}
}
