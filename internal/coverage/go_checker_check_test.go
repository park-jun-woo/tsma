package coverage

import "testing"

func TestGoCheckerCheck_InvalidProject(t *testing.T) {
	checker := &GoChecker{}
	_, err := checker.Check("/nonexistent", "fake_test.go", nil)
	if err == nil {
		t.Error("expected error for nonexistent project")
	}
}
