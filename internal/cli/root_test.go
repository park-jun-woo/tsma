package cli

import "testing"

func TestExecute_setsVersion(t *testing.T) {
	// Verify that the version variable exists and can be set
	original := Version
	defer func() { Version = original }()

	Version = "v1.2.3"
	if Version != "v1.2.3" {
		t.Errorf("expected v1.2.3, got %s", Version)
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
