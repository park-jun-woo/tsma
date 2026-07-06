package index

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCollectSourceFilesWalkBranches complements TestCollectSourceFiles (in
// ts_extract_helpers_test.go) with the walk branches it does not reach: a walk
// callback error being swallowed (nonexistent root) and a non-SkipDir skipDir
// error propagating out of the walk.
func TestCollectSourceFilesWalkBranches(t *testing.T) {
	isTS := func(p string) bool { return strings.HasSuffix(p, ".ts") }

	t.Run("nonexistent root swallows walk error", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "does-not-exist")
		files, err := collectSourceFiles(root, isTS, func(string) error { return nil })
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		if len(files) != 0 {
			t.Errorf("files = %v, want none", files)
		}
	})

	t.Run("skipDir error propagates", func(t *testing.T) {
		root := t.TempDir()
		wantErr := os.ErrPermission
		files, err := collectSourceFiles(root, isTS, func(string) error { return wantErr })
		if err != wantErr {
			t.Fatalf("err = %v, want %v", err, wantErr)
		}
		if len(files) != 0 {
			t.Errorf("files = %v, want none", files)
		}
	})
}
