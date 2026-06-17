package coverage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/match"
)

// TestGoCheckerCheck_overlayError covers GoChecker.Check's overlay-write error
// branch: the match carries a non-empty Overlay but .tsma is occupied by a
// regular file, so serializing the overlay fails before go test runs.
func TestGoCheckerCheck_overlayError(t *testing.T) {
	root := t.TempDir()
	pkgDir := filepath.Join(root, "pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	testFile := filepath.Join(pkgDir, "foo_test.go")
	if err := os.WriteFile(testFile, []byte("package pkg\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Occupy .tsma with a regular file so GoOverlayArgs's MkdirAll(.tsma/test) fails.
	if err := os.WriteFile(filepath.Join(root, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := match.TestMatch{
		Files:     []string{testFile},
		TestFuncs: []string{"TestFoo"},
		Overlay:   map[string]string{filepath.Join(pkgDir, "v_test.go"): testFile},
	}
	checker := &GoChecker{}
	if _, err := checker.Check(root, m, nil); err == nil {
		t.Fatal("expected an overlay-write error from Check")
	}
}
