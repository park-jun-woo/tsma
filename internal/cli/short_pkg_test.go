//ff:test feature=cli
package cli

import (
	"path/filepath"
	"testing"
)

func TestShortPkg_RelativeToRoot(t *testing.T) {
	root := filepath.Join("/project")
	pkgDir := filepath.Join("/project", "internal", "cli")
	if got := shortPkg(root, pkgDir); got != filepath.Join("internal", "cli") {
		t.Errorf("shortPkg = %q, want %q", got, filepath.Join("internal", "cli"))
	}
}

// TestShortPkg_FallbackToAbsolute: when filepath.Rel cannot relate the two paths
// (one absolute, one relative), shortPkg falls back to the package dir as-is.
func TestShortPkg_FallbackToAbsolute(t *testing.T) {
	root := "/abs/root"
	pkgDir := "relative/pkg"
	if got := shortPkg(root, pkgDir); got != pkgDir {
		t.Errorf("shortPkg fallback = %q, want %q (the unmodified pkgDir)", got, pkgDir)
	}
}
