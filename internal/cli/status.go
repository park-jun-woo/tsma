//ff:func feature=cli type=command control=sequence
//ff:what Shows overall progress or delegates to detailed endpoint status
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var statusEndpoint string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show overall progress or detailed status of a specific endpoint",
	Long: `Show overall progress summary (total/done/partial/todo counts and percentages),
or detailed chain-level branch coverage for a specific endpoint with --endpoint.`,
	RunE: runStatus,
}

func init() {
	statusCmd.Flags().StringVar(&statusEndpoint, "endpoint", "", "show detailed status for a specific endpoint")
}

func runStatus(cmd *cobra.Command, args []string) error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	if !session.Exists(root) {
		return fmt.Errorf("no session found — run 'tsma next' first to initialize")
	}
	sess, err := session.Load(root)
	if err != nil {
		return fmt.Errorf("load session: %w", err)
	}

	if statusEndpoint != "" {
		return showEndpointStatus(sess, statusEndpoint)
	}

	return showOverallStatus(sess)
}
