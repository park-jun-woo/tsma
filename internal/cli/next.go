//ff:func feature=cli type=command control=sequence
//ff:what Dispatches `tsma next` to the first-pass watermark scan or the interactive rotating cursor
package cli

import (
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Detect test changes, run tests, measure coverage, and advance",
	Long: `The core command. Finds the current TODO function, detects test file
changes, runs tests, measures coverage, and advances to the next function.

The first pass (after reset) measures every function once with its existing
tests: 100% become PASS, while partials and untested functions stay TODO. A
single partial never blocks the functions behind it from being measured. Once
every function has been measured once, subsequent runs surface the remaining
TODOs one at a time via a rotating cursor.`,
	RunE: runNext,
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

	// A prior run may have left the watermark at the end without flipping the
	// flag (e.g. the last first-pass function was a PASS); reconcile here.
	finishFirstPassIfDone(sess)

	if sess.FirstPassDone {
		return runNextInteractive(root, sess)
	}
	return runNextFirstPass(root, sess)
}
