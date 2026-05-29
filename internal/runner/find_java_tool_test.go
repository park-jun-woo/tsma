package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindJavaToolWrapper(t *testing.T) {
	dir := t.TempDir()
	wrapper := filepath.Join(dir, "mvnw")
	if err := os.WriteFile(wrapper, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	path, err := findJavaTool(dir, "maven")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != wrapper {
		t.Errorf("path = %q, want wrapper %q", path, wrapper)
	}
}

func TestFindJavaToolFromPath(t *testing.T) {
	binDir := t.TempDir()
	fake := filepath.Join(binDir, "gradle")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir)

	path, err := findJavaTool(t.TempDir(), "gradle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != fake {
		t.Errorf("path = %q, want %q", path, fake)
	}
}

func TestFindJavaToolNotFound(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if _, err := findJavaTool(t.TempDir(), "maven"); err == nil {
		t.Fatal("expected error when maven is missing")
	}
}

func TestFindJavaToolUnknown(t *testing.T) {
	if _, err := findJavaTool(t.TempDir(), "ant"); err == nil {
		t.Fatal("expected error for unknown build tool")
	}
}
