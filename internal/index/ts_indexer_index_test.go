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

func TestTSIndexerIndexIgnoresFile(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "keep.ts"), []byte("export function keep() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Ignored FILE -> file-ignore branch returns nil (no SkipDir).
	if err := os.WriteFile(filepath.Join(dir, "ignoreme.ts"), []byte("export function skipped() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".tsmignore"), []byte("ignoreme.ts\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	idx := &TSIndexer{}
	funcs, err := idx.Index(dir)
	if err != nil {
		t.Fatalf("TSIndexer.Index: %v", err)
	}
	if len(funcs) != 1 || funcs[0].Name != "keep" {
		t.Fatalf("expected only keep, got %+v", funcs)
	}
}

func TestTSIndexerIndexWalkError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; unreadable-dir walk error does not apply")
	}
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "keep.ts"), []byte("export function keep() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "locked")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "inner.ts"), []byte("export function inner() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(sub, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(sub, 0o755) })

	idx := &TSIndexer{}
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
