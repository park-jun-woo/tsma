//ff:func feature=cli type=command control=sequence
//ff:what Runs tests per package and updates function statuses based on pass or fail
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/runner"
	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run tests per package and update function statuses",
	Long: `Check re-matches test files, groups functions by package, and runs
tests per package. Pass → DONE, fail → FAIL with output, no test → TODO.`,
	RunE: runCheck,
}

func runCheck(cmd *cobra.Command, args []string) error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	sess, err := loadOrAnalyze(root)
	if err != nil {
		return err
	}

	fmt.Printf("Checking %d functions...\n", len(sess.Functions))

	m := match.NewMatcher(sess.Lang)
	groups := groupByPackage(root, sess, m)

	r := runner.NewRunner(sess.Lang)
	checkAllGroups(root, r, groups)

	sess.RecalcSummary()
	if err := session.Save(root, sess); err != nil {
		return fmt.Errorf("save session: %w", err)
	}

	fmt.Printf("\n%d functions — DONE: %d | FAIL: %d | TODO: %d\n",
		sess.Summary.Total, sess.Summary.Done, sess.Summary.Fail, sess.Summary.Todo)

	return nil
}
