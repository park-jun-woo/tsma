//ff:func feature=cli type=command control=sequence
//ff:what Measures branch coverage for a single function or all DONE functions
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var coverCmd = &cobra.Command{
	Use:   "cover [func]",
	Short: "Measure branch coverage for a function or all DONE functions",
	Long: `Without arguments, measure coverage for all DONE functions grouped by package.
With a function name, measure coverage for that specific function.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCover,
}

func runCover(cmd *cobra.Command, args []string) error {
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

	checker := coverage.NewChecker(sess.Lang)

	if len(args) == 1 {
		return coverSingleFunction(root, sess, checker, args[0])
	}
	return coverAllDone(root, sess, checker)
}
