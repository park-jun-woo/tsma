//ff:func feature=cli type=command control=sequence
//ff:what Shows overall progress summary or detailed status for a specific function
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var statusFunc string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show overall progress or detailed status of a specific function",
	Long: `Show overall progress summary (testable/done/partial/todo counts and percentages),
or detailed coverage info for a specific function with --func.`,
	RunE: runStatus,
}

func init() {
	statusCmd.Flags().StringVar(&statusFunc, "func", "", "show detailed status for a specific function")
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

	if statusFunc != "" {
		fn := sess.FindFunction(statusFunc)
		if fn == nil {
			return fmt.Errorf("function not found: %s", statusFunc)
		}
		printFuncDetail(fn)
		return nil
	}

	// Overall summary.
	s := sess.Summary
	testable := float64(s.Testable)
	if testable == 0 {
		fmt.Println("No testable functions found.")
		return nil
	}

	fmt.Printf("%d testable functions\n", s.Testable)
	fmt.Printf("DONE:    %3d (%.1f%%)\n", s.Done, float64(s.Done)/testable*100)
	fmt.Printf("PARTIAL: %3d (%.1f%%)\n", s.Partial, float64(s.Partial)/testable*100)
	fmt.Printf("TODO:    %3d (%.1f%%)\n", s.Todo, float64(s.Todo)/testable*100)

	g := sess.Graph
	fmt.Printf("\nCall graph: %d nodes, %d edges, %d entry points, %d dead\n",
		g.Nodes, g.Edges, g.EntryPoints, g.Dead)

	return nil
}
