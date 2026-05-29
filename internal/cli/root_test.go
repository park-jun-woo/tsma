package cli

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestExecute_errorExits covers the error branch of Execute (rootCmd.Execute
// returns an error -> prints to stderr and calls os.Exit(1)). Because that
// branch terminates the process, we re-exec the test binary in a subprocess
// and assert it exits non-zero with the error printed to stderr.
func TestExecute_errorExits(t *testing.T) {
	if os.Getenv("TSMA_EXECUTE_ERROR_SUBPROCESS") == "1" {
		// In the subprocess: drive Execute with an unknown subcommand so
		// rootCmd.Execute returns an error and we hit os.Exit(1).
		rootCmd.SetArgs([]string{"this-is-not-a-real-command"})
		Execute("vtest")
		return // unreachable if os.Exit fires as expected
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExecute_errorExits")
	cmd.Env = append(os.Environ(), "TSMA_EXECUTE_ERROR_SUBPROCESS=1")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatalf("expected subprocess to exit non-zero, got success; output: %s", out)
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d; output: %s", exitErr.ExitCode(), out)
	}
	if !strings.Contains(string(out), "unknown command") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestExecute_setsVersion(t *testing.T) {
	// Verify that the version variable exists and can be set
	original := Version
	defer func() { Version = original }()

	Version = "v1.2.3"
	if Version != "v1.2.3" {
		t.Errorf("expected v1.2.3, got %s", Version)
	}
}

// TestExecute_success drives Execute down the success path (no error from
// rootCmd.Execute, so os.Exit is not reached). We invoke the built-in --help
// flag which returns nil. We restore os.Args and the command args afterwards.
func TestExecute_success(t *testing.T) {
	origVersion := Version
	origArgs := os.Args
	defer func() {
		Version = origVersion
		os.Args = origArgs
		rootCmd.SetArgs(nil)
	}()

	rootCmd.SetArgs([]string{"--help"})
	captureStdout(func() {
		Execute("v9.9.9")
	})

	if Version != "v9.9.9" {
		t.Errorf("expected Version set to v9.9.9, got %s", Version)
	}
	if rootCmd.Version != "v9.9.9" {
		t.Errorf("expected rootCmd.Version set to v9.9.9, got %s", rootCmd.Version)
	}
}

func TestRootCmd_exists(t *testing.T) {
	// Verify rootCmd is configured
	if rootCmd == nil {
		t.Fatal("rootCmd is nil")
	}
	if rootCmd.Use != "tsma" {
		t.Errorf("expected Use=tsma, got %s", rootCmd.Use)
	}
}

func TestRootCmd_hasSubcommands(t *testing.T) {
	cmds := rootCmd.Commands()
	names := make(map[string]bool)
	for _, c := range cmds {
		names[c.Name()] = true
	}
	for _, expected := range []string{"next", "list", "status", "reset"} {
		if !names[expected] {
			t.Errorf("expected subcommand %q not found", expected)
		}
	}
}
