//ff:func feature=cli type=command control=iteration dimension=1
//ff:what Runs the 'list' command showing all endpoints with pagination
package cli

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var (
	listPage int
	listSize int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all endpoints with their progress status",
	Long: `List all endpoints with their current status (DONE/PARTIAL/TODO).
Supports pagination with --page and --size flags.`,
	RunE: runList,
}

func init() {
	listCmd.Flags().IntVar(&listPage, "page", 1, "page number (1-based)")
	listCmd.Flags().IntVar(&listSize, "size", 20, "items per page")
}

func runList(cmd *cobra.Command, args []string) error {
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

	// Print summary header.
	fmt.Printf("%d endpoints — DONE: %d | PARTIAL: %d | TODO: %d\n\n",
		sess.Summary.Total, sess.Summary.Done, sess.Summary.Partial, sess.Summary.Todo)

	// Calculate pagination.
	total := len(sess.Endpoints)
	if listPage < 1 {
		listPage = 1
	}
	if listSize < 1 {
		listSize = 20
	}
	start := (listPage - 1) * listSize
	if start >= total {
		fmt.Println("(no endpoints on this page)")
		return nil
	}
	end := start + listSize
	if end > total {
		end = total
	}

	// Find the longest endpoint name for alignment.
	maxName := 0
	for _, ep := range sess.Endpoints[start:end] {
		if len(ep.Name) > maxName {
			maxName = len(ep.Name)
		}
	}

	// Print endpoints.
	for _, ep := range sess.Endpoints[start:end] {
		status := strings.ToUpper(ep.Status)
		extra := ""
		if ep.Status != "partial" || len(ep.UncoveredBranches) == 0 {
			fmt.Printf("  %-*s  %s%s\n", maxName, ep.Name, status, extra)
			continue
		}
		lines := make([]string, 0, len(ep.UncoveredBranches))
		for _, l := range ep.UncoveredBranches {
			lines = append(lines, fmt.Sprintf("%d", l))
		}
		extra = fmt.Sprintf(" (uncovered: line %s)", strings.Join(lines, ", "))
		fmt.Printf("  %-*s  %s%s\n", maxName, ep.Name, status, extra)
	}

	// Print pagination info.
	totalPages := (total + listSize - 1) / listSize
	if totalPages > 1 {
		fmt.Printf("\nPage %d/%d (use --page N to navigate)\n", listPage, totalPages)
	}

	return nil
}
