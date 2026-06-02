//ff:func feature=cli type=command control=sequence
//ff:what `tsma rescan`: reset 없이 소스를 재스캔해 함수 집합만 갱신(진행상태 보존)
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/session"
	"github.com/spf13/cobra"
)

var rescanCmd = &cobra.Command{
	Use:   "rescan",
	Short: "Re-scan source and sync the function set without resetting progress",
	Long: `Re-scan the source tree and reconcile the session's function set with it:
newly added functions become TODO (and are measured), removed functions are
dropped, and existing functions keep their PASS/DONE/TODO progress. Unlike
'reset --all', rescan preserves progress — use it after refactors that add,
extract, move, or delete functions.`,
	RunE: runRescan,
}

func runRescan(cmd *cobra.Command, args []string) error {
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

	added, removed, err := reconcileSession(root, sess, true)
	if err != nil {
		return err
	}
	if err := session.Save(root, sess); err != nil {
		return fmt.Errorf("save session: %w", err)
	}

	fmt.Printf("Rescan: %d functions (+%d new, -%d removed)\n", len(sess.Functions), added, removed)
	return nil
}
