//ff:func feature=cli type=command control=sequence
//ff:what Lists all functions with pagination, status filtering, and sorting
package cli

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var (
	listPage   int
	listSize   int
	listStatus string
	listSort   string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all functions with their progress status",
	Long: `List all functions with their current status (DONE/PARTIAL/TODO).
Supports pagination, status filtering, and sorting.`,
	RunE: runList,
}

func init() {
	listCmd.Flags().IntVar(&listPage, "page", 1, "page number (1-based)")
	listCmd.Flags().IntVar(&listSize, "size", 20, "items per page")
	listCmd.Flags().StringVar(&listStatus, "status", "", "filter by status: todo, partial, done, dead")
	listCmd.Flags().StringVar(&listSort, "sort", "priority", "sort by: priority, name, file")
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

	// Filter functions.
	statusFilter := strings.ToLower(listStatus)
	filtered, err := filterFunctions(sess.Functions, statusFilter)
	if err != nil {
		return err
	}

	// Sort functions.
	if err := sortFunctions(filtered, strings.ToLower(listSort)); err != nil {
		return err
	}

	// Print summary header.
	if statusFilter == "dead" {
		fmt.Printf("%d dead functions\n\n", len(filtered))
	} else {
		fmt.Printf("%d functions — DONE: %d | PARTIAL: %d | TODO: %d\n\n",
			sess.Summary.Testable, sess.Summary.Done, sess.Summary.Partial, sess.Summary.Todo)
	}

	// Calculate pagination.
	total := len(filtered)
	if listPage < 1 {
		listPage = 1
	}
	if listSize < 1 {
		listSize = 20
	}
	start := (listPage - 1) * listSize
	if start >= total {
		fmt.Println("(no functions on this page)")
		return nil
	}
	end := start + listSize
	if end > total {
		end = total
	}

	page := filtered[start:end]
	maxName := maxFuncNameLen(page)
	printFuncList(page, maxName)

	// Print pagination info.
	totalPages := (total + listSize - 1) / listSize
	if totalPages > 1 {
		fmt.Printf("\nPage %d/%d (use --page N to navigate)\n", listPage, totalPages)
	}

	return nil
}
