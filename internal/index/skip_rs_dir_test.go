package index

import (
	"path/filepath"
	"testing"
)

func TestSkipRsDir(t *testing.T) {
	skipped := []string{"target", ".git", ".tsma"}
	for _, d := range skipped {
		if got := skipRsDir(filepath.Join("/proj", d)); got != filepath.SkipDir {
			t.Errorf("skipRsDir(%q) = %v, want SkipDir", d, got)
		}
	}

	if got := skipRsDir("/proj/src"); got != nil {
		t.Errorf("skipRsDir(src) = %v, want nil", got)
	}
}
