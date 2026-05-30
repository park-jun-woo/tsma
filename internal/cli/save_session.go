//ff:func feature=cli type=helper control=sequence
//ff:what Persists the session, wrapping any save error with context
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/session"
)

// saveSession persists sess and wraps any error, shared by the next modes.
func saveSession(root string, sess *model.Session) error {
	if err := session.Save(root, sess); err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}
