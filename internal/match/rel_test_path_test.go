package match

import (
	"path/filepath"
	"testing"
)

func TestRelTestPath(t *testing.T) {
	root := filepath.Join("home", "proj")
	abs := filepath.Join(root, "internal", "gen", "x_test.go")
	got := relTestPath(root, abs)
	want := filepath.Join("internal", "gen", "x_test.go")
	if got != want {
		t.Errorf("relTestPath = %q, want %q", got, want)
	}
}

func TestRelTestPathSameDir(t *testing.T) {
	root := filepath.Join("home", "proj")
	abs := filepath.Join(root, "x_test.go")
	if got := relTestPath(root, abs); got != "x_test.go" {
		t.Errorf("relTestPath = %q, want x_test.go", got)
	}
}

func TestRelTestPathNeverAbsolute(t *testing.T) {
	// Mixing an absolute target with a relative root makes filepath.Rel fail;
	// the fallback must return the base name, never an absolute path.
	got := relTestPath("relative/root", filepath.Join(string(filepath.Separator), "abs", "y_test.go"))
	if filepath.IsAbs(got) {
		t.Errorf("relTestPath returned absolute path %q", got)
	}
	if got != "y_test.go" {
		t.Errorf("relTestPath fallback = %q, want base y_test.go", got)
	}
}
