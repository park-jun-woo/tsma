//ff:func feature=cli type=helper control=sequence
//ff:what Loads existing session or auto-analyzes the project if none exists
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/session"
)

// loadOrAnalyze loads an existing session or performs initial analysis.
func loadOrAnalyze(root string) (*model.Session, error) {
	if session.Exists(root) {
		sess, err := session.Load(root)
		if err != nil {
			return nil, fmt.Errorf("load session: %w", err)
		}
		return sess, nil
	}

	fmt.Fprintln(os.Stderr, "No session found. Analyzing project...")
	sess, err := analyzeProject(root)
	if err != nil {
		return nil, fmt.Errorf("analyze project: %w", err)
	}
	if err := session.Save(root, sess); err != nil {
		return nil, fmt.Errorf("save session: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Session created: %d functions indexed (%s)\n\n",
		len(sess.Functions), sess.Lang)

	return sess, nil
}
