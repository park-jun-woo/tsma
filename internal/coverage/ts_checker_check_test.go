package coverage

import "testing"

func TestTSCheckerCheck_InvalidProject(t *testing.T) {
	checker := &TSChecker{}
	_, err := checker.Check("/nonexistent", "fake_test.ts", nil)
	if err == nil {
		t.Error("expected error for nonexistent project")
	}
}
