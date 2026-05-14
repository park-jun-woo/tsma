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
