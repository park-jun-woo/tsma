//ff:func feature=cli type=command control=sequence
//ff:what Lists all functions with pagination
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var listPage int

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all functions with their progress status",
	Long:  `List all functions with their current status (DONE/PARTIAL/TODO) and coverage.`,
	RunE:  runList,
}

func init() {
	listCmd.Flags().IntVar(&listPage, "page", 1, "page number (1-based)")
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

	functions := sess.Functions

	// Print summary header.
	fmt.Printf("%d functions — DONE: %d | PARTIAL: %d | TODO: %d\n\n",
		sess.Summary.Total, sess.Summary.Done, sess.Summary.Partial, sess.Summary.Todo)

	// Calculate pagination.
	const pageSize = 20
	total := len(functions)
	if listPage < 1 {
		listPage = 1
	}
	start := (listPage - 1) * pageSize
	if start >= total {
		fmt.Println("(no functions on this page)")
		return nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	page := functions[start:end]
	maxName := maxFuncNameLen(page)
	printFuncList(page, maxName)

	// Print pagination info.
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages > 1 {
		fmt.Printf("\nPage %d/%d (use --page N to navigate)\n", listPage, totalPages)
	}

	return nil
}
