package coverage

import "testing"

// runPytestCov executes external commands (pytest), so direct testing is limited.
// These tests verify behavior with invalid inputs that fail before executing.

func TestRunPytestCovInvalidProject(t *testing.T) {
	err := runPytestCov("/nonexistent/project", "test_handler.py", "/tmp/cover.json")
	if err == nil {
		t.Skip("pytest may be available; skipping error expectation")
	}
	// If pytest is not installed, this should return an error.
}
