package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindReadmeFrom_inStartDir(t *testing.T) {
	dir := t.TempDir()
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# x"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := findReadmeFrom(dir)
	if got != readme {
		t.Errorf("expected %q, got %q", readme, got)
	}
}

func TestFindReadmeFrom_walksUpToParent(t *testing.T) {
	parent := t.TempDir()
	readme := filepath.Join(parent, "README.md")
	if err := os.WriteFile(readme, []byte("# x"), 0o644); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(parent, "a", "b")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}

	got := findReadmeFrom(child)
	if got != readme {
		t.Errorf("expected parent README %q, got %q", readme, got)
	}
}

func TestFindReadmeFrom_fallbackToGitHub(t *testing.T) {
	// A temp dir has no README.md anywhere up to the filesystem root.
	dir := t.TempDir()

	got := findReadmeFrom(dir)
	if got != readmeGitHubURL {
		t.Errorf("expected GitHub URL %q, got %q", readmeGitHubURL, got)
	}
}

func TestFindReadmeFrom_ignoresDirectoryNamedReadme(t *testing.T) {
	dir := t.TempDir()
	// A directory (not a file) named README.md must not be treated as the README.
	if err := os.Mkdir(filepath.Join(dir, "README.md"), 0o755); err != nil {
		t.Fatal(err)
	}

	got := findReadmeFrom(dir)
	if got != readmeGitHubURL {
		t.Errorf("expected GitHub URL (dir README.md ignored), got %q", got)
	}
}
