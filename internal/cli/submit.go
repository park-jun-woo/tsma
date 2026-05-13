//ff:func feature=cli type=command control=sequence
//ff:what Submits a test file for a function, validates it passes, and checks branch coverage
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/runner"
	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit <func-name> <test-file>",
	Short: "Submit a test file for a function",
	Long: `Submit a test file for verification against a function.
Two-stage validation:
  1. Pass check — run the test against the original code
  2. Coverage check — verify branch coverage for the function

Results:
  DONE    — test passes and 100% branch coverage
  PARTIAL — test passes but coverage is below 100% (retry up to 2 times)`,
	Args: cobra.ExactArgs(2),
	RunE: runSubmit,
}

func runSubmit(cmd *cobra.Command, args []string) error {
	funcName := args[0]
	testFile := args[1]

	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	// Check test file exists.
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		return fmt.Errorf("test file not found: %s", testFile)
	}

	// Load session.
	if !session.Exists(root) {
		return fmt.Errorf("no session found — run 'tsma next' first to initialize")
	}
	sess, err := session.Load(root)
	if err != nil {
		return fmt.Errorf("load session: %w", err)
	}

	// Find the function.
	fn := sess.FindFunction(funcName)
	if fn == nil {
		return fmt.Errorf("function not found: %s", funcName)
	}

	// Stage 1: Run tests.
	fmt.Printf("[1/2] Running tests... ")
	r := runner.NewRunner(sess.Lang)
	result, err := r.Run(root, testFile)
	if err != nil {
		return fmt.Errorf("run tests: %w", err)
	}
	if !result.Pass {
		fmt.Println("FAIL")
		fmt.Println(result.Output)
		fmt.Println("NOT accepted.")
		return nil
	}
	fmt.Println("PASS")

	// Stage 2: Check coverage.
	fmt.Printf("[2/2] Checking branch coverage for %s...\n", funcName)
	checker := coverage.NewChecker(sess.Lang)
	report, err := checker.Check(root, testFile, fn)
	if err != nil {
		return fmt.Errorf("coverage check: %w", err)
	}

	printCoverageReport(report)

	// Copy test file and update function.
	relTest, err := session.CopyTestFile(root, testFile)
	if err != nil {
		return fmt.Errorf("copy test file: %w", err)
	}
	fn.TestFile = relTest

	prevPct := fn.CoveragePct
	applySubmitResult(fn, report, prevPct)
	sess.RecalcSummary()
	printSubmitOutcome(fn, funcName, report, sess)

	if err := session.Save(root, sess); err != nil {
		return fmt.Errorf("save session: %w", err)
	}

	return nil
}
