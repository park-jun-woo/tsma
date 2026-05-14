package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCoverProfileValid(t *testing.T) {
	dir := t.TempDir()
	content := `mode: set
github.com/example/pkg/handler.go:10.2,20.5 3 1
github.com/example/pkg/handler.go:22.2,30.5 2 0
`
	path := filepath.Join(dir, "cover.out")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	blocks, err := parseCoverProfile(path)
	if err != nil {
		t.Fatalf("parseCoverProfile: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks, want 2", len(blocks))
	}
	if blocks[0].count != 1 {
		t.Errorf("blocks[0].count = %d, want 1", blocks[0].count)
	}
	if blocks[1].count != 0 {
		t.Errorf("blocks[1].count = %d, want 0", blocks[1].count)
	}
}

func TestParseCoverProfileFileNotFound(t *testing.T) {
	_, err := parseCoverProfile("/nonexistent/cover.out")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestParseCoverProfileSkipsModeAndInvalid(t *testing.T) {
	dir := t.TempDir()
	content := `mode: atomic
invalid line without proper format
github.com/example/pkg/handler.go:10.2,20.5 3 1
`
	path := filepath.Join(dir, "cover.out")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	blocks, err := parseCoverProfile(path)
	if err != nil {
		t.Fatalf("parseCoverProfile: %v", err)
	}
	if len(blocks) != 1 {
		t.Fatalf("got %d blocks, want 1 (mode and invalid line skipped)", len(blocks))
	}
}

func TestParseCoverProfileEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cover.out")
	if err := os.WriteFile(path, []byte("mode: set\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	blocks, err := parseCoverProfile(path)
	if err != nil {
		t.Fatalf("parseCoverProfile: %v", err)
	}
	if len(blocks) != 0 {
		t.Fatalf("got %d blocks, want 0", len(blocks))
	}
}
