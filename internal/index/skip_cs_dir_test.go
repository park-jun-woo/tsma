package index

import (
	"path/filepath"
	"testing"
)

func TestSkipCsDir(t *testing.T) {
	for _, base := range []string{"bin", "obj", ".vs", ".git", ".tsma"} {
		if err := skipCsDir(filepath.Join("proj", base)); err != filepath.SkipDir {
			t.Errorf("skipCsDir(%q) = %v, want SkipDir", base, err)
		}
	}
	for _, base := range []string{"src", "Models", "Services"} {
		if err := skipCsDir(filepath.Join("proj", base)); err != nil {
			t.Errorf("skipCsDir(%q) = %v, want nil", base, err)
		}
	}
}
