//ff:func feature=cli type=command control=sequence
//ff:what Deletes the entire session when called with --all flag
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var resetAll bool

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Delete the entire session",
	Long: `Delete the entire session with --all flag.
After reset, the next 'tsma next' will re-analyze the project.`,
	RunE: runReset,
}

func init() {
	resetCmd.Flags().BoolVar(&resetAll, "all", false, "delete the entire session")
}

func runReset(cmd *cobra.Command, args []string) error {
	if !resetAll {
		return fmt.Errorf("use --all to delete the entire session")
	}

	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	if err := session.Delete(root); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	fmt.Println("Session deleted. Next `tsma next` will re-analyze the project.")
	return nil
}
