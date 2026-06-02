//ff:func feature=cli type=command control=sequence
//ff:what Shows overall progress summary with counts, percentages, and average coverage
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show overall progress summary",
	Long:  `Show overall progress summary (total/pass/done/todo counts, percentages, and average coverage).`,
	RunE:  runStatus,
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

	// Reflect source changes so status never reports stale counts (BUG-004).
	// Structure-only (measure=false): new functions surface as TODO without
	// running tests; measurement is left to `tsma next`/`tsma rescan`. Always
	// save: even with no add/remove, positional metadata and CheckedAt may have
	// been refreshed.
	added, removed, err := reconcileSession(root, sess, false)
	if err != nil {
		return err
	}
	if err := session.Save(root, sess); err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	if added > 0 || removed > 0 {
		fmt.Fprintf(os.Stderr, "source changed: +%d new, -%d removed (run 'tsma next')\n", added, removed)
	}

	s := sess.Summary
	total := float64(s.Total)
	if total == 0 {
		fmt.Println("No functions found.")
		return nil
	}

	fmt.Printf("%d functions\n", s.Total)
	fmt.Printf("PASS: %4d (%.1f%%)\n", s.Pass, float64(s.Pass)/total*100)
	fmt.Printf("DONE: %4d (%.1f%%)\n", s.Done, float64(s.Done)/total*100)
	fmt.Printf("TODO: %4d (%.1f%%)\n", s.Todo, float64(s.Todo)/total*100)

	avg := computeAverageCoverage(sess)
	fmt.Printf("\nCoverage average: %.0f%%\n", avg)

	return nil
}
