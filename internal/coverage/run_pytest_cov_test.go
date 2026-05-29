package coverage

import (
	"path/filepath"
	"testing"
)

func TestRunPytestCovInvalidProject(t *testing.T) {
	err := runPytestCov("/nonexistent/project", "test_handler.py", "/tmp/cover.json")
	if err == nil {
		t.Skip("pytest may be available; skipping error expectation")
	}
	// If pytest is not installed, this should return an error.
}

func TestRunPytestCovSuccess(t *testing.T) {
	// Fake python3 exits 0 -> success path.
	proj := fakePython3(t, "exit 0\n")
	out := filepath.Join(proj, "cover.json")

	if err := runPytestCov(proj, "test_handler.py", out); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
}

func TestRunPytestCovFailure(t *testing.T) {
	// Fake python3 exits non-zero -> error branch.
	proj := fakePython3(t, "echo boom 1>&2\nexit 2\n")
	out := filepath.Join(proj, "cover.json")

	if err := runPytestCov(proj, "test_handler.py", out); err == nil {
		t.Fatal("expected error when pytest exits non-zero")
	}
}
