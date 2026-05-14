package coverage

import "testing"

func TestPyCheckerCheck_InvalidProject(t *testing.T) {
	checker := &PyChecker{}
	_, err := checker.Check("/nonexistent", "fake_test.py", nil)
	if err == nil {
		t.Error("expected error for nonexistent project")
	}
}
