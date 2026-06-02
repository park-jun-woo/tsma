//ff:func feature=cli type=command control=sequence
//ff:what Runs `tsma next`: analyze-on-first-run (batch-measure all), then the interactive rotating cursor
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Detect test changes, run tests, measure coverage, and advance",
	Long: `The core command. On the first run (no session) it indexes every
function, matches test files, and batch-measures coverage for all of them in a
single pass: 100%% functions become PASS, partials become measured TODOs, and
untested functions stay TODO. Subsequent runs surface the remaining TODOs one at
a time via a rotating cursor, re-measuring a TODO whose test changed and never
stalling on a single partial.`,
	RunE: runNext,
}

func init() {
	nextCmd.Flags().Int("max-attempts", defaultMaxAttempts,
		"auto-DONE a partial function after this many tsma next presentations (must be >= 1)")
}

func runNext(cmd *cobra.Command, args []string) error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	sess, fresh, err := loadOrAnalyze(root)
	if err != nil {
		return err
	}

	// A loaded session may be stale: functions added/extracted/removed since the
	// last index are invisible until reconciled, which is how a fully-PASS stale
	// session falsely reported "All functions complete!" (BUG-004). Re-scan the
	// source and merge the current function set, preserving progress. A
	// freshly-analyzed session already reflects current source, so skip it there.
	if !fresh {
		added, removed, rErr := reconcileSession(root, sess, true)
		if rErr != nil {
			return rErr
		}
		if added > 0 || removed > 0 {
			fmt.Fprintf(os.Stderr, "Source changed: +%d new, -%d removed\n", added, removed)
		}
	}

	if err := resolveMaxAttempts(cmd, sess); err != nil {
		return err
	}

	// After analysis/reconcile the session is fully measured, so `tsma next` is
	// the interactive rotating cursor over the remaining TODOs.
	return runNextInteractive(root, sess)
}
