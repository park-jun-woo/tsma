package treesitter

import (
	"errors"
	"os/exec"
	"testing"
)

// lookOrSkip resolves a coreutil needed by the retry tests, skipping when the
// host lacks it (keeps the suite green on minimal images).
func lookOrSkip(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s not available: %v", name, err)
	}
	return path
}

func TestRunOnceCapturesStdout(t *testing.T) {
	echo := lookOrSkip(t, "echo")
	out, err := runOnce(echo, []string{"hello"})
	if err != nil {
		t.Fatalf("runOnce error: %v", err)
	}
	if string(out) != "hello\n" {
		t.Errorf("stdout = %q, want %q", out, "hello\n")
	}
}

func TestRunOnceStartFailure(t *testing.T) {
	out, err := runOnce("tsma-no-such-binary-xyz", nil)
	if err == nil {
		t.Fatal("expected error for missing binary")
	}
	if len(out) != 0 {
		t.Errorf("stdout = %q, want empty", out)
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		t.Error("missing binary should be a start error, not ExitError")
	}
}

func TestShouldRetryRun(t *testing.T) {
	falseCmd := lookOrSkip(t, "false")
	_, exitErr := runOnce(falseCmd, nil) // real *exec.ExitError, empty stdout
	_, startErr := runOnce("tsma-no-such-binary-xyz", nil)

	cases := []struct {
		name string
		out  []byte
		err  error
		want bool
	}{
		{"success-nil-err", []byte("<xml/>"), nil, false},
		{"nonempty-with-exit", []byte("<xml/>"), exitErr, false},
		{"empty-exit-error", nil, exitErr, true},
		{"empty-start-error", nil, startErr, false},
		{"empty-no-error", nil, nil, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := shouldRetryRun(c.out, c.err); got != c.want {
				t.Errorf("shouldRetryRun = %v, want %v", got, c.want)
			}
		})
	}
}

func TestRunReturnsOutputFirstTry(t *testing.T) {
	echo := lookOrSkip(t, "echo")
	out, err := Run(echo, "", []string{"x"})
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected non-empty stdout")
	}
}

func TestRunGivesUpAfterRetries(t *testing.T) {
	falseCmd := lookOrSkip(t, "false")
	old := treeSitterRetryBackoff
	treeSitterRetryBackoff = 0 // no sleeps in test
	defer func() { treeSitterRetryBackoff = old }()

	out, err := Run(falseCmd, "g", []string{"x"})
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if out != nil {
		t.Errorf("out = %q, want nil", out)
	}
}
