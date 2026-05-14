package coverage

import "testing"

func TestFindCoveragePython(t *testing.T) {
	result := findCoveragePython()
	// Should return either "python3" or "python"
	if result != "python3" && result != "python" {
		t.Errorf("findCoveragePython() = %q, want 'python3' or 'python'", result)
	}
}
