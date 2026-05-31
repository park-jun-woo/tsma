//ff:func feature=cli type=helper control=sequence
//ff:what Returns the session's auto-DONE threshold, falling back to the default when unset
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// effectiveMaxAttempts returns the session's configured auto-DONE threshold,
// falling back to defaultMaxAttempts when it is unset (0). resolveMaxAttempts
// normally persists a value during `tsma next`, but reading through this helper
// keeps the threshold correct for any caller (including sessions constructed
// without going through the flag-resolution path) so an unset 0 never collapses
// the threshold into "auto-DONE immediately".
func effectiveMaxAttempts(sess *model.Session) int {
	if sess.MaxAttempts < 1 {
		return defaultMaxAttempts
	}
	return sess.MaxAttempts
}
