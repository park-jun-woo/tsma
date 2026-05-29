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

func TestPyIndexerIndexIgnoresFile(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "keep.py"), []byte("def keep(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A specific ignored FILE -> file-ignore branch returns nil (no SkipDir).
	if err := os.WriteFile(filepath.Join(dir, "ignoreme.py"), []byte("def skipped(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".tsmignore"), []byte("ignoreme.py\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}
	if len(funcs) != 1 || funcs[0].Name != "keep" {
		t.Fatalf("expected only keep, got %+v", funcs)
	}
}

func TestPyIndexerIndexWalkError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; unreadable-dir walk error does not apply")
	}
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "keep.py"), []byte("def keep(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "locked")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "inner.py"), []byte("def inner(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(sub, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(sub, 0o755) })

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("expected success despite walk error, got: %v", err)
	}
	found := false
	for _, fn := range funcs {
		if fn.Name == "keep" {
			found = true
		}
	}
	if !found {
		t.Error("expected keep to still be indexed")
	}
}

func TestPyIndexerIndexRespectsIgnore(t *testing.T) {
	dir := t.TempDir()

	// Create a Python file at root
	if err := os.WriteFile(filepath.Join(dir, "main.py"), []byte("def keep(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a Python file inside "excluded/" directory
	excludedDir := filepath.Join(dir, "excluded")
	if err := os.MkdirAll(excludedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(excludedDir, "skip.py"), []byte("def skip(): pass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create .tsmignore excluding "excluded/"
	if err := os.WriteFile(filepath.Join(dir, ".tsmignore"), []byte("excluded/\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &PyIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("PyIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (excluded/ skipped), got %d", len(funcs))
	}
	if funcs[0].Name != "keep" {
		t.Errorf("expected keep, got %s", funcs[0].Name)
	}
}
