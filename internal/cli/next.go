//ff:func feature=cli type=command control=sequence
//ff:what Shows the next TODO function to test
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
	Short: "Show the next incomplete function to test",
	Long: `Show the next incomplete function (TODO or PARTIAL) in order.
If no session exists, automatically analyze the project first.`,
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
		fmt.Fprintf(os.Stderr, "Session created: %d functions indexed (%s)\n\n",
			len(sess.Functions), sess.Lang)
	}

	// Find the first TODO or PARTIAL function.
	next := findFirstIncomplete(sess)

	if next == nil {
		fmt.Println("All functions are DONE!")
		return nil
	}

	statusLabel := "TODO"
	if next.Status == model.StatusPartial {
		statusLabel = "PARTIAL"
	}
	fmt.Printf("%s  %s\n", next.Name, statusLabel)
	fmt.Printf("  file: %s:%d-%d\n", next.File, next.StartLine, next.EndLine)

	return nil
}
