package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectGo(t *testing.T) {
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

func TestDetectPythonPyproject(t *testing.T) {
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

func TestDetectPythonRequirements(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "requirements.txt"), []byte("flask==2.0\n"), 0o644); err != nil {
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

func TestDetectPythonSetupPy(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "setup.py"), []byte("from setuptools import setup\n"), 0o644); err != nil {
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

func TestDetectGoPriority(t *testing.T) {
	// When both go.mod and package.json exist, Go wins.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "go" {
		t.Errorf("Lang = %q, want go (go.mod takes priority)", lf.Lang)
	}
}
