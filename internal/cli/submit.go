//ff:func feature=cli type=command control=iteration dimension=1
//ff:what Submits a test file, validates it passes, and checks branch coverage
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/runner"
	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit <endpoint> <test-file>",
	Short: "Submit a test file for an endpoint",
	Long: `Submit a test file for verification against an endpoint.
Two-stage validation:
  1. Pass check — run the test against the original code
  2. Coverage check — verify branch coverage for the endpoint's function chain

Results:
  DONE    — test passes and 100% branch coverage
  PARTIAL — test passes but coverage is below 100% (uncovered branches shown)
  FAIL    — test does not pass (rejected)`,
	Args: cobra.ExactArgs(2),
	RunE: runSubmit,
}

func runSubmit(cmd *cobra.Command, args []string) error {
	endpointName := args[0]
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

	// Find the endpoint.
	ep := sess.FindEndpoint(endpointName)
	if ep == nil {
		return fmt.Errorf("endpoint not found: %s", endpointName)
	}

	// Stage 1: Run tests.
	fmt.Printf("[1/2] Running tests against original code... ")
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
	fmt.Printf("[2/2] Checking branch coverage for %s...\n", endpointName)
	checker := coverage.NewChecker(sess.Lang)
	report, err := checker.Check(root, testFile, ep)
	if err != nil {
		return fmt.Errorf("coverage check: %w", err)
	}

	// Print coverage details.
	for _, fc := range report.Funcs {
		pct := fmt.Sprintf("%.0f%%", fc.CoveredPct)
		blocks := fmt.Sprintf("(%d/%d blocks)", fc.CoveredBlocks, fc.TotalBlocks)
		fmt.Printf("  %-40s %s %s\n", fc.File+":"+fmt.Sprintf("%d-%d", fc.StartLine, fc.EndLine), pct, blocks)
		for _, line := range fc.UncoveredLines {
			fmt.Printf("    UNCOVERED: branch at line %d\n", line)
		}
	}

	// Copy test file and update session.
	relTest, err := session.CopyTestFile(root, testFile)
	if err != nil {
		return fmt.Errorf("copy test file: %w", err)
	}
	ep.TestFile = relTest
	ep.Coverage = make(map[string]float64)
	for _, fc := range report.Funcs {
		ep.Coverage[fc.Key] = fc.CoveredPct
	}

	if report.AllCovered {
		ep.Status = model.StatusDone
		ep.UncoveredBranches = nil
		sess.RecalcSummary()
		remaining := sess.Summary.Todo + sess.Summary.Partial
		fmt.Printf("DONE — %s complete (%d remaining)\n", endpointName, remaining)
	} else {
		ep.Status = model.StatusPartial
		ep.UncoveredBranches = nil
		for _, ub := range report.Uncovered {
			ep.UncoveredBranches = append(ep.UncoveredBranches, ub.Line)
		}
		sess.RecalcSummary()
		uncoveredCount := len(report.Uncovered)
		fmt.Printf("PARTIAL — %s needs %d more branch(es) covered\n", endpointName, uncoveredCount)
	}

	if err := session.Save(root, sess); err != nil {
		return fmt.Errorf("save session: %w", err)
	}

	return nil
}
