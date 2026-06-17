package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/tsma/internal/match"
)

// TestGoOverlayArgs_EmptyReturnsNil covers the empty-overlay guard: the manual
// submit path keeps its vet-on behavior (no flags, no file written).
func TestGoOverlayArgs_EmptyReturnsNil(t *testing.T) {
	args, err := GoOverlayArgs(t.TempDir(), nil)
	if err != nil {
		t.Fatalf("empty overlay: %v", err)
	}
	if args != nil {
		t.Fatalf("empty overlay must return nil args, got %v", args)
	}
}

// TestGoOverlayArgs_WritesAndReturnsFlags covers the success path: the overlay map
// is serialized under .tsma/test/overlay.json and the activating flags are returned.
func TestGoOverlayArgs_WritesAndReturnsFlags(t *testing.T) {
	root := t.TempDir()
	overlay := map[string]string{"/v/zzz_test.go": "/b/gen.go"}
	args, err := GoOverlayArgs(root, overlay)
	if err != nil {
		t.Fatalf("GoOverlayArgs: %v", err)
	}
	jsonPath := filepath.Join(root, ".tsma", "test", "overlay.json")
	want := []string{"-overlay", jsonPath, "-vet=off"}
	if len(args) != 3 || args[0] != want[0] || args[1] != want[1] || args[2] != want[2] {
		t.Fatalf("args = %v, want %v", args, want)
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("overlay.json must be written: %v", err)
	}
	var got struct{ Replace map[string]string }
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("overlay.json must be valid JSON: %v", err)
	}
	if got.Replace["/v/zzz_test.go"] != "/b/gen.go" {
		t.Fatalf("overlay mapping not persisted: %+v", got.Replace)
	}
}

// TestGoOverlayArgs_MkdirError covers the directory-creation failure: .tsma is a
// regular file, so MkdirAll(.tsma/test) cannot succeed.
func TestGoOverlayArgs_MkdirError(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := GoOverlayArgs(root, map[string]string{"a": "b"}); err == nil {
		t.Fatal("expected an error when .tsma cannot be a directory")
	}
}

// TestGoOverlayArgs_WriteError covers the file-write failure: the overlay.json
// path is occupied by a directory, so WriteFile fails.
func TestGoOverlayArgs_WriteError(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".tsma", "test", "overlay.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := GoOverlayArgs(root, map[string]string{"a": "b"}); err == nil {
		t.Fatal("expected an error when overlay.json cannot be written")
	}
}

// TestGoRunnerRun_overlayError covers GoRunner.Run's overlay-write error branch:
// a valid module with a non-empty Overlay, but .tsma occupied by a file so the
// overlay serialization fails before the test command is built.
func TestGoRunnerRun_overlayError(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/o\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "foo_test.go"), []byte("package o\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".tsma"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	m := match.TestMatch{
		Files:     []string{filepath.Join(dir, "foo_test.go")},
		TestFuncs: []string{"TestFoo"},
		Overlay:   map[string]string{filepath.Join(dir, "v_test.go"): filepath.Join(dir, "foo_test.go")},
	}
	r := &GoRunner{}
	if _, err := r.Run(dir, m); err == nil {
		t.Fatal("expected an overlay-write error from Run")
	}
}
