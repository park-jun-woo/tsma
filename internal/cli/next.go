//ff:func feature=cli type=command control=sequence
//ff:what Shows the next incomplete function prioritized by incoming edges
package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show the next incomplete function to test",
	Long: `Show the next incomplete function (TODO or PARTIAL) prioritized by incoming
edges (callers). If no session exists, automatically analyze the project first
(language detection, function indexing, call graph building).`,
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

	// Collect incomplete functions (exclude dead code).
	candidates := collectIncomplete(sess.Functions)

	if len(candidates) == 0 {
		fmt.Println("All functions are DONE!")
		return nil
	}

	// Sort by priority: incoming edges desc, leaf functions first, then file path.
	sort.Slice(candidates, func(i, j int) bool {
		ci := len(candidates[i].Callers)
		cj := len(candidates[j].Callers)
		if ci != cj {
			return ci > cj
		}
		li := len(candidates[i].Callees) == 0
		lj := len(candidates[j].Callees) == 0
		if li != lj {
			return li
		}
		if candidates[i].File != candidates[j].File {
			return candidates[i].File < candidates[j].File
		}
		return candidates[i].StartLine < candidates[j].StartLine
	})

	next := &candidates[0]

	// Print function info.
	statusLabel := "TODO"
	if next.Status == model.StatusPartial {
		statusLabel = "PARTIAL"
	}
	fmt.Printf("%s  %s  (priority: %d callers)\n", next.Name, statusLabel, len(next.Callers))
	fmt.Printf("  file: %s:%d-%d\n", next.File, next.StartLine, next.EndLine)

	if len(next.Callers) > 0 {
		printNextCallers(next.Callers, 3, sess)
	}

	if len(next.Callees) > 0 {
		printNextCallees(next.Callees, sess)
	}

	return nil
}
