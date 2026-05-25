package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTSIndexerIndexBasic(t *testing.T) {
	dir := t.TempDir()

	content := `export async function handleLogin(req: Request) {
  return await authService.login(req.body);
}

function helperFunc() {
  return "ok";
}
`
	if err := os.WriteFile(filepath.Join(dir, "handler.ts"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}

	if len(funcs) < 2 {
		t.Fatalf("expected at least 2 functions, got %d", len(funcs))
	}
}

func TestTSIndexerIndexSkipsTestFiles(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "handler.ts"), []byte("export function handler() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "handler.test.ts"), []byte("describe('handler', () => {});\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "handler.spec.ts"), []byte("describe('handler', () => {});\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (test/spec excluded), got %d", len(funcs))
	}
}

func TestTSIndexerIndexSkipsNodeModules(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "index.ts"), []byte("export function main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	nmDir := filepath.Join(dir, "node_modules", "lib")
	if err := os.MkdirAll(nmDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nmDir, "index.ts"), []byte("export function libFunc() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (node_modules excluded), got %d", len(funcs))
	}
}

func TestTSIndexerIndexEmptyProject(t *testing.T) {
	dir := t.TempDir()

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("expected 0 functions, got %d", len(funcs))
	}
}

func TestTSIndexerIndexRespectsIgnore(t *testing.T) {
	dir := t.TempDir()

	// Create a TS file at root
	if err := os.WriteFile(filepath.Join(dir, "index.ts"), []byte("export function keep() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a TS file inside "generated/" directory
	genDir := filepath.Join(dir, "generated")
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(genDir, "auto.ts"), []byte("export function auto() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create .tsmignore excluding "generated/"
	if err := os.WriteFile(filepath.Join(dir, ".tsmignore"), []byte("generated/\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}

	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (generated/ skipped), got %d", len(funcs))
	}
	if funcs[0].Name != "keep" {
		t.Errorf("expected keep, got %s", funcs[0].Name)
	}
}
