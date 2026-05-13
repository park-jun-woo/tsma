//ff:func feature=cli type=command control=sequence
//ff:what Configures the root cobra command and runs the CLI
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

// rootCmd is the top-level command.
var rootCmd = &cobra.Command{
	Use:   "tsma",
	Short: "Extract tests from legacy code with LLM agents",
	Long: `TestMaster manages endpoint-level test extraction from legacy codebases.
It detects endpoints, traces function call chains, validates submitted tests,
and tracks branch coverage progress.`,
	Version: Version,
}

// Execute runs the root command.
func Execute(version string) {
	Version = version
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(submitCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(resetCmd)
}
