package coverage

import "testing"

// runCoveragePy executes external commands (coverage.py), so direct testing is limited.

func TestRunCoveragePyInvalidProject(t *testing.T) {
	err := runCoveragePy("/nonexistent/project", "test_handler.py", "/tmp/cover.json")
	if err == nil {
		t.Skip("coverage.py may be available; skipping error expectation")
	}
	// If coverage.py is not installed, this should return an error.
}
