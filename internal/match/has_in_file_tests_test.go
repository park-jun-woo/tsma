package match

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHasInFileTests(t *testing.T) {
	dir := t.TempDir()

	withTests := filepath.Join(dir, "a.rs")
	if err := os.WriteFile(withTests, []byte("pub fn f() {}\n#[cfg(test)]\nmod tests {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !hasInFileTests(withTests) {
		t.Error("expected hasInFileTests true for file with #[cfg(test)]")
	}

	without := filepath.Join(dir, "b.rs")
	if err := os.WriteFile(without, []byte("pub fn f() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if hasInFileTests(without) {
		t.Error("expected hasInFileTests false for file without #[cfg(test)]")
	}

	if hasInFileTests(filepath.Join(dir, "missing.rs")) {
		t.Error("expected hasInFileTests false for missing file")
	}
}
