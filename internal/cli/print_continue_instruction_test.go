package cli

import (
	"strings"
	"testing"
)

func TestPrintContinueInstruction(t *testing.T) {
	output := captureStdout(func() {
		printContinueInstruction()
	})
	if !strings.Contains(output, "tsma next") {
		t.Errorf("expected output to mention 'tsma next', got: %q", output)
	}
	if !strings.Contains(output, "next function") {
		t.Errorf("expected output to mention 'next function', got: %q", output)
	}
}
