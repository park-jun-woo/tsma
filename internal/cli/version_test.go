package cli

import (
	"strings"
	"testing"
)

func TestRunVersion_printsVersionAndReadme(t *testing.T) {
	orig := Version
	defer func() { Version = orig }()
	Version = "v9.9.9"

	out := captureStdout(func() {
		if err := runVersion(versionCmd, nil); err != nil {
			t.Fatalf("runVersion returned error: %v", err)
		}
	})

	if !strings.Contains(out, "tsma version v9.9.9") {
		t.Errorf("expected version line, got: %q", out)
	}
	if !strings.Contains(out, "Must read ") {
		t.Errorf("expected README pointer, got: %q", out)
	}
	if !strings.Contains(out, "README.md") {
		t.Errorf("expected README.md reference, got: %q", out)
	}
}
