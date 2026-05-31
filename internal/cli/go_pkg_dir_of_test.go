package cli

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/match"
)

// TestGoPkgDirOf covers goPkgDirOf directly: an absolute first file is returned
// as its own directory unchanged, while a relative first file is resolved
// against the project root before taking its directory.
func TestGoPkgDirOf(t *testing.T) {
	root := filepath.Join(string(filepath.Separator)+"proj", "root")

	t.Run("relative path resolved against root", func(t *testing.T) {
		m := match.TestMatch{Files: []string{filepath.Join("internal", "cli", "x_test.go")}}
		want := filepath.Join(root, "internal", "cli")
		if got := goPkgDirOf(root, m); got != want {
			t.Fatalf("goPkgDirOf(root, relative) = %q, want %q", got, want)
		}
	})

	t.Run("relative root-level file resolves to root", func(t *testing.T) {
		m := match.TestMatch{Files: []string{"main_test.go"}}
		if got := goPkgDirOf(root, m); got != root {
			t.Fatalf("goPkgDirOf(root, root-level) = %q, want %q", got, root)
		}
	})

	t.Run("absolute path used verbatim", func(t *testing.T) {
		abs := filepath.Join(string(filepath.Separator)+"abs", "pkg", "y_test.go")
		m := match.TestMatch{Files: []string{abs}}
		want := filepath.Dir(abs)
		if got := goPkgDirOf(root, m); got != want {
			t.Fatalf("goPkgDirOf(root, absolute) = %q, want %q", got, want)
		}
	})
}
