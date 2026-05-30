package cli

import (
	"strings"
	"testing"
)

func TestPrintAllComplete(t *testing.T) {
	out := captureStdout(func() { printAllComplete() })
	if !strings.Contains(out, "All functions complete!") {
		t.Errorf("expected completion banner, got %q", out)
	}
}
