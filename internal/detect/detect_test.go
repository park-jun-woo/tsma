package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectGo(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example\n\ngo 1.22\n\nrequire github.com/gin-gonic/gin v1.10.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "go" {
		t.Errorf("Lang = %q, want go", lf.Lang)
	}
	if lf.Framework != "gin" {
		t.Errorf("Framework = %q, want gin", lf.Framework)
	}
}

func TestDetectGoUnknownFramework(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "go" {
		t.Errorf("Lang = %q, want go", lf.Lang)
	}
	if lf.Framework != "unknown" {
		t.Errorf("Framework = %q, want unknown", lf.Framework)
	}
}

func TestDetectUnsupported(t *testing.T) {
	dir := t.TempDir()
	_, err := Detect(dir)
	if err == nil {
		t.Error("expected error for unsupported language")
	}
}

func TestDetectTypeScript(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "typescript" {
		t.Errorf("Lang = %q, want typescript", lf.Lang)
	}
}

func TestDetectPython(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte("[tool.poetry]"), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "python" {
		t.Errorf("Lang = %q, want python", lf.Lang)
	}
}

// ---------------------------------------------------------------------------
// TypeScript framework detection
// ---------------------------------------------------------------------------

func TestDetectTSExpress(t *testing.T) {
	dir := t.TempDir()
	pkg := `{
  "name": "my-app",
  "dependencies": {
    "express": "^4.18.0"
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "typescript" {
		t.Errorf("Lang = %q, want typescript", lf.Lang)
	}
	if lf.Framework != "express" {
		t.Errorf("Framework = %q, want express", lf.Framework)
	}
}

func TestDetectTSNextjs(t *testing.T) {
	dir := t.TempDir()
	pkg := `{
  "name": "my-app",
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.0.0"
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "typescript" {
		t.Errorf("Lang = %q, want typescript", lf.Lang)
	}
	if lf.Framework != "nextjs" {
		t.Errorf("Framework = %q, want nextjs", lf.Framework)
	}
}

func TestDetectTSUnknown(t *testing.T) {
	dir := t.TempDir()
	pkg := `{
  "name": "my-app",
  "dependencies": {
    "lodash": "^4.17.0"
  }
}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "typescript" {
		t.Errorf("Lang = %q, want typescript", lf.Lang)
	}
	if lf.Framework != "unknown" {
		t.Errorf("Framework = %q, want unknown", lf.Framework)
	}
}

// ---------------------------------------------------------------------------
// Python framework detection
// ---------------------------------------------------------------------------

func TestDetectPyFastapi(t *testing.T) {
	dir := t.TempDir()
	pyproject := `[project]
dependencies = [
    "fastapi>=0.100.0",
    "uvicorn>=0.20.0",
]
`
	if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte(pyproject), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "python" {
		t.Errorf("Lang = %q, want python", lf.Lang)
	}
	if lf.Framework != "fastapi" {
		t.Errorf("Framework = %q, want fastapi", lf.Framework)
	}
}

func TestDetectPyDjango(t *testing.T) {
	dir := t.TempDir()
	pyproject := `[project]
dependencies = [
    "django>=4.2",
]
`
	if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte(pyproject), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "python" {
		t.Errorf("Lang = %q, want python", lf.Lang)
	}
	if lf.Framework != "django" {
		t.Errorf("Framework = %q, want django", lf.Framework)
	}
}

func TestDetectPyFromRequirements(t *testing.T) {
	dir := t.TempDir()
	req := `fastapi>=0.100.0
uvicorn>=0.20.0
pydantic>=2.0
`
	if err := os.WriteFile(filepath.Join(dir, "requirements.txt"), []byte(req), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "python" {
		t.Errorf("Lang = %q, want python", lf.Lang)
	}
	if lf.Framework != "fastapi" {
		t.Errorf("Framework = %q, want fastapi", lf.Framework)
	}
}

func TestDetectPyFromSetupPy(t *testing.T) {
	dir := t.TempDir()
	setup := `from setuptools import setup

setup(
    name="my-app",
    install_requires=[
        "django>=4.2",
        "djangorestframework>=3.14",
    ],
)
`
	if err := os.WriteFile(filepath.Join(dir, "setup.py"), []byte(setup), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "python" {
		t.Errorf("Lang = %q, want python", lf.Lang)
	}
	if lf.Framework != "django" {
		t.Errorf("Framework = %q, want django", lf.Framework)
	}
}

// ---------------------------------------------------------------------------
// Go framework detection (echo, chi)
// ---------------------------------------------------------------------------

func TestDetectGoEcho(t *testing.T) {
	dir := t.TempDir()
	gomod := `module example

go 1.22

require github.com/labstack/echo/v4 v4.12.0
`
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "go" {
		t.Errorf("Lang = %q, want go", lf.Lang)
	}
	if lf.Framework != "echo" {
		t.Errorf("Framework = %q, want echo", lf.Framework)
	}
}

func TestDetectGoChi(t *testing.T) {
	dir := t.TempDir()
	gomod := `module example

go 1.22

require github.com/go-chi/chi/v5 v5.0.12
`
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "go" {
		t.Errorf("Lang = %q, want go", lf.Lang)
	}
	if lf.Framework != "chi" {
		t.Errorf("Framework = %q, want chi", lf.Framework)
	}
}

// ---------------------------------------------------------------------------
// Helper function tests (containsImport, indexOf)
// ---------------------------------------------------------------------------

func TestContainsImport(t *testing.T) {
	if !containsImport("require github.com/gin-gonic/gin v1.10.0", "github.com/gin-gonic/gin") {
		t.Error("expected containsImport to find gin")
	}
	if containsImport("require github.com/foo/bar v1.0.0", "github.com/gin-gonic/gin") {
		t.Error("expected containsImport to not find gin")
	}
	if containsImport("", "anything") {
		t.Error("expected containsImport to return false for empty content")
	}
}

func TestIndexOf(t *testing.T) {
	if idx := indexOf("hello world", "world"); idx != 6 {
		t.Errorf("indexOf = %d, want 6", idx)
	}
	if idx := indexOf("hello world", "xyz"); idx != -1 {
		t.Errorf("indexOf = %d, want -1", idx)
	}
	if idx := indexOf("", "a"); idx != -1 {
		t.Errorf("indexOf empty = %d, want -1", idx)
	}
}
