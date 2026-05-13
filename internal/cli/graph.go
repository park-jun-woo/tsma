//ff:func feature=cli type=command control=sequence
//ff:what Queries the call graph for a function's callers and callees or shows dead code
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var (
	graphFunc    string
	graphDead    bool
	graphCallers bool
	graphCallees bool
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Query the call graph",
	Long: `Query the call graph for a function's callers and callees.
Without --func, shows the overall graph summary.
With --dead, lists all dead code functions.
Use --callers or --callees to filter the output.`,
	RunE: runGraph,
}

func init() {
	graphCmd.Flags().StringVar(&graphFunc, "func", "", "show callers/callees for a specific function")
	graphCmd.Flags().BoolVar(&graphDead, "dead", false, "list dead code functions")
	graphCmd.Flags().BoolVar(&graphCallers, "callers", false, "show only callers (requires --func)")
	graphCmd.Flags().BoolVar(&graphCallees, "callees", false, "show only callees (requires --func)")
}

func runGraph(cmd *cobra.Command, args []string) error {
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

	if graphDead {
		printDeadFunctions(sess.Functions)
		return nil
	}

	if graphFunc != "" {
		fn := sess.FindFunction(graphFunc)
		if fn == nil {
			return fmt.Errorf("function not found: %s", graphFunc)
		}
		printFuncGraph(sess, fn)
		return nil
	}

	// Overall graph summary.
	g := sess.Graph
	fmt.Printf("Call graph summary:\n")
	fmt.Printf("  Nodes:        %d\n", g.Nodes)
	fmt.Printf("  Edges:        %d\n", g.Edges)
	fmt.Printf("  Entry points: %d\n", g.EntryPoints)
	fmt.Printf("  Dead code:    %d\n", g.Dead)

	return nil
}
