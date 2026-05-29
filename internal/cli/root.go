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
	Short: "Function-level test status tracker for legacy codebases",
	Long: `TestMaster tracks function-level test status for legacy codebases.
It indexes all functions, matches test files, runs tests to determine
pass/fail status, and measures branch coverage.`,
	Version: Version,
}

// Execute runs the root command.
func Execute(version string) {
	Version = version
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("tsma version %s\nMust read %s\n", version, findReadme()))
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(versionCmd)
}
