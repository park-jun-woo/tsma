package cli

import (
	"os"
	"testing"
)

func TestGetProjectRoot_returnsWorkingDir(t *testing.T) {
	root, err := getProjectRoot()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if root != wd {
		t.Errorf("expected %q, got %q", wd, root)
	}
}

// TestGetProjectRoot_getwdError exercises the error branch by removing the
// working directory out from under the process, which makes os.Getwd fail.
func TestGetProjectRoot_getwdError(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// Restore a valid cwd so other tests are unaffected.
	t.Cleanup(func() { _ = os.Chdir(orig) })

	dir, err := os.MkdirTemp("", "tsma-getwd")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	// Remove the directory we're standing in. On Linux os.Getwd then errors.
	if err := os.Remove(dir); err != nil {
		t.Skipf("could not remove cwd, skipping: %v", err)
	}

	_, err = getProjectRoot()
	if err == nil {
		t.Skip("os.Getwd did not fail on this platform; error branch not exercised")
	}
}
