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
