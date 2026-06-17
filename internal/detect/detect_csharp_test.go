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

// TestDetectCSharpDirectGlobMatch covers the glob-match branch (line 19):
// detectCSharp returns true directly when a *.csproj file is present.
func TestDetectCSharpDirectGlobMatch(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "App.csproj"), []byte("<Project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !detectCSharp(dir) {
		t.Error("expected detectCSharp=true when a .csproj is present")
	}
}

// TestDetectCSharpDirectPropsMatch covers the Directory.Build.props branch
// (line 23): no project/solution glob hit, but the fixed-name props file marks
// the project as C#.
func TestDetectCSharpDirectPropsMatch(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Directory.Build.props"), []byte("<Project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !detectCSharp(dir) {
		t.Error("expected detectCSharp=true when Directory.Build.props is present")
	}
}
