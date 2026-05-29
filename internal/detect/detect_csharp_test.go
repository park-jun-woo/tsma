package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectCSharpCsproj(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "App.csproj"), []byte("<Project/>"), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "csharp" {
		t.Errorf("Lang = %q, want csharp", lf.Lang)
	}
}

func TestDetectCSharpSln(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Solution.sln"), []byte("Microsoft Visual Studio Solution File\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "csharp" {
		t.Errorf("Lang = %q, want csharp", lf.Lang)
	}
}

func TestDetectCSharpDirectoryBuildProps(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Directory.Build.props"), []byte("<Project/>"), 0o644); err != nil {
		t.Fatal(err)
	}

	lf, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if lf.Lang != "csharp" {
		t.Errorf("Lang = %q, want csharp", lf.Lang)
	}
}

func TestDetectCSharpNone(t *testing.T) {
	if detectCSharp(t.TempDir()) {
		t.Error("expected no C# marker in an empty dir")
	}
}
