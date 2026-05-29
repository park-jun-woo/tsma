package cli

import (
	"strings"
	"testing"
)

func TestPrintNextInstruction(t *testing.T) {
	output := captureStdout(func() {
		printNextInstruction()
	})
	if !strings.Contains(output, "tsma next") {
		t.Errorf("expected output to mention 'tsma next', got: %q", output)
	}
	if !strings.Contains(output, "completing the test") {
		t.Errorf("expected output to mention 'completing the test', got: %q", output)
	}
}
