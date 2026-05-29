//ff:func feature=cli type=command control=sequence
//ff:what Defines the version subcommand printing the version and a README pointer
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd prints the version followed by a pointer to the README.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and a pointer to the README",
	RunE:  runVersion,
}

// runVersion prints the version and the resolved README location.
func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("tsma version %s\nMust read %s\n", Version, findReadme())
	return nil
}
