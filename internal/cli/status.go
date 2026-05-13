//ff:func feature=cli type=command control=sequence
//ff:what Shows overall progress summary with counts and percentages
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show overall progress summary",
	Long:  `Show overall progress summary (total/done/fail/todo counts and percentages).`,
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

	s := sess.Summary
	total := float64(s.Total)
	if total == 0 {
		fmt.Println("No functions found.")
		return nil
	}

	fmt.Printf("%d functions\n", s.Total)
	fmt.Printf("DONE: %4d (%.1f%%)\n", s.Done, float64(s.Done)/total*100)
	fmt.Printf("FAIL: %4d (%.1f%%)\n", s.Fail, float64(s.Fail)/total*100)
	fmt.Printf("TODO: %4d (%.1f%%)\n", s.Todo, float64(s.Todo)/total*100)

	return nil
}
