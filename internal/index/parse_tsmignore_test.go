package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTsmIgnoreFileNotFound(t *testing.T) {
	patterns := ParseTsmIgnore("/nonexistent/.tsmignore")
	if patterns != nil {
		t.Errorf("expected nil for nonexistent file, got %v", patterns)
	}
}

func TestParseTsmIgnoreEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tsmignore")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	patterns := ParseTsmIgnore(path)
	if len(patterns) != 0 {
		t.Errorf("expected empty slice, got %v", patterns)
	}
}

func TestParseTsmIgnoreSkipsCommentsAndBlanks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tsmignore")
	content := `# comment line

vendor/
  # indented comment

*.log
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	patterns := ParseTsmIgnore(path)
	if len(patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d: %v", len(patterns), patterns)
	}
	if patterns[0] != "vendor/" {
		t.Errorf("patterns[0] = %q, want %q", patterns[0], "vendor/")
	}
	if patterns[1] != "*.log" {
		t.Errorf("patterns[1] = %q, want %q", patterns[1], "*.log")
	}
}

func TestParseTsmIgnoreNormalPatterns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tsmignore")
	content := `node_modules/
dist/
*.generated.go
internal/tmp/*.go
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	patterns := ParseTsmIgnore(path)
	expected := []string{"node_modules/", "dist/", "*.generated.go", "internal/tmp/*.go"}
	if len(patterns) != len(expected) {
		t.Fatalf("expected %d patterns, got %d: %v", len(expected), len(patterns), patterns)
	}
	for i, want := range expected {
		if patterns[i] != want {
			t.Errorf("patterns[%d] = %q, want %q", i, patterns[i], want)
		}
	}
}
