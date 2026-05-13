//ff:func feature=cli type=command control=sequence
//ff:what Resets a function to TODO or deletes the entire session
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var resetAll bool

var resetCmd = &cobra.Command{
	Use:   "reset [func-name]",
	Short: "Reset a function to TODO or delete the entire session",
	Long: `Reset a specific function to TODO status, or use --all to delete the entire
session. After --all, the next 'tsma next' will re-analyze the project.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runReset,
}

func init() {
	resetCmd.Flags().BoolVar(&resetAll, "all", false, "delete the entire session")
}

func runReset(cmd *cobra.Command, args []string) error {
	root, err := getProjectRoot()
	if err != nil {
		return err
	}

	if resetAll {
		if err := session.Delete(root); err != nil {
			return fmt.Errorf("delete session: %w", err)
		}
		fmt.Println("Session deleted. Next `tsma next` will re-analyze the project.")
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("specify a function name or use --all")
	}

	funcName := args[0]

	if !session.Exists(root) {
		return fmt.Errorf("no session found — nothing to reset")
	}
	sess, err := session.Load(root)
	if err != nil {
		return fmt.Errorf("load session: %w", err)
	}

	fn := sess.FindFunction(funcName)
	if fn == nil {
		return fmt.Errorf("function not found: %s", funcName)
	}

	fn.Status = model.StatusTodo
	fn.TestFile = ""
	fn.FailOutput = ""
	sess.RecalcSummary()

	if err := session.Save(root, sess); err != nil {
		return fmt.Errorf("save session: %w", err)
	}

	remaining := sess.Summary.Todo + sess.Summary.Fail
	fmt.Printf("%s reset to TODO (%d remaining)\n", funcName, remaining)
	return nil
}
