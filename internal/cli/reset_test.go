package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunReset_withoutAllFlag(t *testing.T) {
	resetAll = false
	err := runReset(nil, nil)
	if err == nil {
		t.Fatal("expected error when --all is not set")
	}
}

// TestRunReset_getProjectRootError covers the getProjectRoot error branch
// (line 32) by removing the cwd so os.Getwd() fails. --all must be set so we
// get past the earlier guard.
func TestRunReset_getProjectRootError(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(orig)

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.Remove(dir); err != nil {
		t.Skipf("could not remove cwd: %v", err)
	}
	if _, gErr := os.Getwd(); gErr == nil {
		t.Skip("os.Getwd did not fail after removing cwd on this platform")
	}

	resetAll = true
	if err := runReset(nil, nil); err == nil {
		t.Fatal("expected error when getProjectRoot fails")
	}
}

func TestRunReset_withAllFlag(t *testing.T) {
	dir := t.TempDir()
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.WriteFile(filepath.Join(sessDir, "session.json"), []byte("{}"), 0o644)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	resetAll = true
	output := captureStdout(func() {
		err := runReset(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	// Verify session directory was deleted
	if _, err := os.Stat(sessDir); err == nil {
		t.Error("expected .tsma directory to be deleted")
	}

	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestRunReset_deleteError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; read-only parent permission check does not apply")
	}
	dir := t.TempDir()

	// Create the .tsma dir with a file inside, then make the project root
	// read-only so RemoveAll of .tsma fails (cannot unlink within a
	// non-writable parent directory).
	sessDir := filepath.Join(dir, ".tsma")
	os.MkdirAll(sessDir, 0o755)
	os.WriteFile(filepath.Join(sessDir, "session.json"), []byte("{}"), 0o644)
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	resetAll = true
	err := runReset(nil, nil)
	if err == nil {
		t.Skip("RemoveAll did not fail on this platform; delete error branch not exercised")
	}
}

func TestRunReset_noSessionDir(t *testing.T) {
	dir := t.TempDir()

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	resetAll = true
	// Should not error even when no session exists
	captureStdout(func() {
		err := runReset(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
