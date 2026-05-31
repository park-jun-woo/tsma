//ff:func feature=cli type=helper control=sequence
//ff:what Resolves the auto-DONE attempt threshold from the flag, session, or default and stores it
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/spf13/cobra"
)

// defaultMaxAttempts is the fallback presentation threshold for auto-DONE.
const defaultMaxAttempts = 3

// resolveMaxAttempts reconciles the --max-attempts flag, the stored session
// value, and the default (3), persisting the effective value onto the session so
// repeated `tsma next` calls stay consistent. An explicit flag wins and must be
// >= 1 (rejected otherwise); without a flag (or a nil command, as in tests that
// invoke runNext directly) the stored value is kept, and an unset (0) session
// value falls back to the default.
func resolveMaxAttempts(cmd *cobra.Command, sess *model.Session) error {
	if cmd != nil && cmd.Flags().Changed("max-attempts") {
		n, err := cmd.Flags().GetInt("max-attempts")
		if err != nil {
			return err
		}
		if n < 1 {
			return fmt.Errorf("--max-attempts must be >= 1")
		}
		sess.MaxAttempts = n
		return nil
	}
	if sess.MaxAttempts == 0 {
		sess.MaxAttempts = defaultMaxAttempts
	}
	return nil
}
