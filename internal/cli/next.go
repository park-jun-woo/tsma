//ff:func feature=cli type=command control=selection
//ff:what Orchestrates the next-function workflow: detect, run, measure, advance
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Detect test changes, run tests, measure coverage, and advance",
	Long: `The core command. Finds the current TODO function, detects test file
changes, runs tests, measures coverage, and advances to the next function.`,
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

	fn := advanceToNext(sess)
	if fn == nil {
		fmt.Println("All functions complete!")
		if err := session.Save(root, sess); err != nil {
			return fmt.Errorf("save session: %w", err)
		}
		return nil
	}

	changed, testFile := detectTestChange(root, sess.Lang, fn)
	if testFile == "" {
		printTodoFunction(fn, "")
		printNextInstruction()
		if err := session.Save(root, sess); err != nil {
			return fmt.Errorf("save session: %w", err)
		}
		return nil
	}

	if !changed {
		printTodoFunction(fn, testFile)
		printNextInstruction()
		if err := session.Save(root, sess); err != nil {
			return fmt.Errorf("save session: %w", err)
		}
		return nil
	}

	result := runAndMeasure(root, sess.Lang, fn, testFile)

	switch result.outcome {
	case outcomeTestFail:
		fn.TestFile = testFile
		fn.TestMtime = result.mtime
		fn.FailOutput = result.failOutput
		fmt.Fprintf(os.Stderr, "FAIL  %s\n", fn.Name)
		fmt.Fprintf(os.Stderr, "  %s\n", result.failOutput)
		printTodoFunction(fn, testFile)
		printNextInstruction()

	case outcomePass:
		fn.Status = model.StatusPass
		fn.TestFile = testFile
		fn.TestMtime = result.mtime
		fn.CoveragePct = 100
		fn.FailOutput = ""
		sess.CurrentIndex++
		fmt.Printf("PASS  %s  100%%\n", fn.Name)
		next := advanceToNext(sess)
		if next != nil {
			printContinueInstruction()
			printTodoFunction(next, next.TestFile)
			printNextInstruction()
		} else {
			fmt.Println("All functions complete!")
		}

	case outcomeDone:
		fn.Status = model.StatusDone
		fn.TestFile = testFile
		fn.TestMtime = result.mtime
		fn.CoveragePct = result.coveragePct
		fn.FailOutput = ""
		sess.CurrentIndex++
		fmt.Printf("DONE  %s  %.0f%%\n", fn.Name, result.coveragePct)
		next := advanceToNext(sess)
		if next != nil {
			printContinueInstruction()
			printTodoFunction(next, next.TestFile)
			printNextInstruction()
		} else {
			fmt.Println("All functions complete!")
		}

	case outcomeRetry:
		fn.TestFile = testFile
		fn.TestMtime = result.mtime
		fn.Attempt = result.attempt
		fn.CoveragePct = result.coveragePct
		fn.FailOutput = ""
		fmt.Printf("PARTIAL  %s  %.0f%%  (attempt %d — improve coverage)\n",
			fn.Name, result.coveragePct, result.attempt)
		printUncovered(result.uncovered)
		printTodoFunction(fn, testFile)
		printNextInstruction()
	}

	sess.RecalcSummary()
	if err := session.Save(root, sess); err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}
