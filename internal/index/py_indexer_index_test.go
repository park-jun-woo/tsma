package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPyIndexerIndexBasic(t *testing.T) {
	dir := t.TempDir()

	content := `def handle_login(request):
    return auth.login(request)

async def fetch_data(url):
    return await get(url)
`
	if err := os.WriteFile(filepath.Join(dir, "handler.py"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}

	if len(funcs) != 2 {
		t.Fatalf("expected 2 functions, got %d", len(funcs))
	}
}

func TestPyIndexerIndexSkipsTestFiles(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "handler.py"), []byte("def handler(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "test_handler.py"), []byte("def test_handler(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (test excluded), got %d", len(funcs))
	}
}

func TestPyIndexerIndexSkipsExcludedDirs(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "main.py"), []byte("def main(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cacheDir := filepath.Join(dir, "__pycache__")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "cached.py"), []byte("def cached(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (__pycache__ excluded), got %d", len(funcs))
	}
}

func TestPyIndexerIndexEmptyProject(t *testing.T) {
	dir := t.TempDir()

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("expected 0 functions, got %d", len(funcs))
	}
}
