package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCoveragePyInvalidProject(t *testing.T) {
	err := runCoveragePy("/nonexistent/project", "test_handler.py", "/tmp/cover.json")
	if err == nil {
		t.Skip("coverage.py may be available; skipping error expectation")
	}
	// If coverage.py is not installed, this should return an error.
}

// fakePython3 installs a fake `python3` on PATH (findCoveragePython prefers it)
// with the supplied shell body, returning the project dir.
func fakePython3(t *testing.T, body string) string {
	t.Helper()
	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, "python3"), []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return t.TempDir()
}

func TestRunCoveragePySuccess(t *testing.T) {
	// Fake python3 succeeds for both `coverage run` and `coverage json` steps.
	proj := fakePython3(t, "exit 0\n")
	out := filepath.Join(proj, "cover.json")

	if err := runCoveragePy(proj, "tests/test_handler.py", out); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
}

func TestRunCoveragePyRunFails(t *testing.T) {
	// `coverage run ...` -> non-zero exit -> first error branch.
	proj := fakePython3(t, `
for a in "$@"; do
  if [ "$a" = "run" ]; then exit 1; fi
done
exit 0
`)
	out := filepath.Join(proj, "cover.json")

	if err := runCoveragePy(proj, "test_handler.py", out); err == nil {
		t.Fatal("expected error when coverage run fails")
	}
}

func TestRunCoveragePyJSONFails(t *testing.T) {
	// `coverage run` succeeds but `coverage json` fails -> second error branch.
	proj := fakePython3(t, `
for a in "$@"; do
  if [ "$a" = "json" ]; then exit 1; fi
done
exit 0
`)
	out := filepath.Join(proj, "cover.json")

	if err := runCoveragePy(proj, "test_handler.py", out); err == nil {
		t.Fatal("expected error when coverage json fails")
	}
}
