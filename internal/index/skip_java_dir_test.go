package index

import (
	"path/filepath"
	"testing"
)

func TestSkipJavaDir(t *testing.T) {
	skipped := []string{"target", "build", ".gradle", ".git", ".tsma"}
	for _, d := range skipped {
		if err := skipJavaDir(filepath.Join("proj", d)); err != filepath.SkipDir {
			t.Errorf("skipJavaDir(%q) = %v, want SkipDir", d, err)
		}
	}
	if err := skipJavaDir(filepath.Join("proj", "src")); err != nil {
		t.Errorf("skipJavaDir(src) = %v, want nil", err)
	}
}
