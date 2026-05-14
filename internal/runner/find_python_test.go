package runner

import "testing"

func TestFindPythonReturnsValidBinary(t *testing.T) {
	result := findPython()
	if result != "python3" && result != "python" {
		t.Errorf("findPython() = %q, want \"python3\" or \"python\"", result)
	}
}
