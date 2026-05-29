package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindReadme_fromWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# x"), 0o644); err != nil {
		t.Fatal(err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	got := findReadme()
	// macOS/Linux may symlink temp dirs; compare resolved paths.
	wantInfo, _ := os.Stat(readme)
	gotInfo, statErr := os.Stat(got)
	if statErr != nil || !os.SameFile(wantInfo, gotInfo) {
		t.Errorf("expected README at %q, got %q", readme, got)
	}
}

func TestFindReadme_fallsBackWhenCwdUnavailable(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig)

	// Enter a directory and then remove it so that os.Getwd fails.
	gone := filepath.Join(t.TempDir(), "gone")
	if err := os.Mkdir(gone, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(gone); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(gone); err != nil {
		t.Skipf("could not remove cwd to force Getwd failure: %v", err)
	}

	if got := findReadme(); got != readmeGitHubURL {
		t.Errorf("expected GitHub URL fallback when cwd unavailable, got %q", got)
	}
}
