package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindDotnetFound(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "dotnet")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	path, err := findDotnet()
	if err != nil {
		t.Fatalf("expected dotnet found, got error: %v", err)
	}
	if path != fake {
		t.Errorf("expected %q, got %q", fake, path)
	}
}

func TestFindDotnetNotFound(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if _, err := findDotnet(); err == nil {
		t.Fatal("expected error when dotnet not on PATH")
	}
}
