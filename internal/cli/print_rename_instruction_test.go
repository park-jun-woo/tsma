package cli

import (
	"strings"
	"testing"
)

func TestPrintRenameInstruction(t *testing.T) {
	output := captureStdout(func() {
		printRenameInstruction("pkg/test_foo_test.go", "pkg/foo_test.go")
	})
	if !strings.Contains(output, "pkg/test_foo_test.go") {
		t.Errorf("expected output to contain the misnamed path, got: %q", output)
	}
	if !strings.Contains(output, "pkg/foo_test.go") {
		t.Errorf("expected output to contain the canonical path, got: %q", output)
	}
	if !strings.Contains(output, "rename") {
		t.Errorf("expected output to mention 'rename', got: %q", output)
	}
	if !strings.Contains(output, "tsma next") {
		t.Errorf("expected output to mention 'tsma next', got: %q", output)
	}
}
