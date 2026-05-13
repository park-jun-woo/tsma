//ff:func feature=cli type=command control=sequence
//ff:what Shows the next TODO or FAIL function to work on
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show the next TODO or FAIL function to work on",
	Long: `Show the next incomplete function (TODO or FAIL) in order.
If no session exists, automatically analyze the project first.`,
	RunE: runNext,
}

func runNext(cmd *cobra.Command, args []string) error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	sess, err := loadOrAnalyze(root)
	if err != nil {
		return err
	}

	next := findFirstIncomplete(sess)
	if next == nil {
		fmt.Println("All functions DONE!")
		return nil
	}

	printNextFunction(next)
	return nil
}
