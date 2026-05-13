//ff:func feature=cli type=command control=iteration dimension=1
//ff:what Shows the next incomplete endpoint with its function call chain
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
	Short: "Show the next incomplete endpoint with its function call chain",
	Long: `Show the next incomplete endpoint (TODO or PARTIAL) along with its function
call chain. If no session exists, automatically analyze the project first
(language detection, endpoint discovery, AST chain tracing).`,
	RunE: runNext,
}

func runNext(cmd *cobra.Command, args []string) error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	var sess *model.Session

	if session.Exists(root) {
		sess, err = session.Load(root)
		if err != nil {
			return fmt.Errorf("load session: %w", err)
		}
	} else {
		// Auto-analyze the project.
		fmt.Fprintln(os.Stderr, "No session found. Analyzing project...")
		sess, err = analyzeProject(root)
		if err != nil {
			return fmt.Errorf("analyze project: %w", err)
		}
		if err := session.Save(root, sess); err != nil {
			return fmt.Errorf("save session: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Session created: %d endpoints detected (%s/%s)\n\n",
			len(sess.Endpoints), sess.Lang, sess.Framework)
	}

	// Find the next TODO or PARTIAL endpoint.
	var next *model.Endpoint
	for i := range sess.Endpoints {
		ep := &sess.Endpoints[i]
		if ep.Status != model.StatusTodo && ep.Status != model.StatusPartial {
			continue
		}
		next = ep
		break
	}

	if next == nil {
		fmt.Println("All endpoints are DONE!")
		return nil
	}

	// Print endpoint info.
	statusLabel := "TODO"
	if next.Status == model.StatusPartial {
		statusLabel = "PARTIAL"
	}
	fmt.Printf("%s\t%s\n", next.Name, statusLabel)
	if next.Method != "" || next.Path != "" {
		fmt.Printf("%s %s\n", next.Method, next.Path)
	}
	fmt.Printf("handler: %s:%d-%d\n", next.Handler.File, next.Handler.StartLine, next.Handler.EndLine)

	if len(next.Chain) == 0 {
		return nil
	}

	fmt.Println("chain:")
	for _, ce := range next.Chain {
		if ce.File != "" {
			fmt.Printf("  -> %-30s %s:%d-%d\n", ce.Func+"()", ce.File, ce.StartLine, ce.EndLine)
			continue
		}
		if ce.Boundary != "" {
			fmt.Printf("  -> %-30s (%s)\n", ce.Func+"()", ce.Boundary)
		}
	}

	return nil
}
