//ff:func feature=cli type=command control=sequence
//ff:what Runs `tsma next`: analyze-on-first-run (batch-measure all), then the interactive rotating cursor
package cli

import (
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

	sess, err := loadOrAnalyze(root)
	if err != nil {
		return err
	}

	if err := resolveMaxAttempts(cmd, sess); err != nil {
		return err
	}

	// The first scan batch-measures every function during analysis (and sets
	// FirstPassDone), so the session is always fully measured here. `tsma next`
	// is therefore purely the interactive rotating cursor over the remaining
	// TODOs — there is no incremental first-pass step anymore.
	return runNextInteractive(root, sess)
}
