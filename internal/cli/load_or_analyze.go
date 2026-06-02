//ff:func feature=cli type=helper control=sequence
//ff:what Loads existing session or auto-analyzes the project if none exists
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/session"
)

// loadOrAnalyze loads an existing session or performs initial analysis. The
// returned fresh flag is true only on the analyze path (a brand-new session
// already reflects current source); callers use it to skip a redundant rescan
// right after a full analysis.
func loadOrAnalyze(root string) (sess *model.Session, fresh bool, err error) {
	if session.Exists(root) {
		sess, err = session.Load(root)
		if err != nil {
			return nil, false, fmt.Errorf("load session: %w", err)
		}
		return sess, false, nil
	}

	fmt.Fprintln(os.Stderr, "No session found. Analyzing project...")
	sess, err = analyzeProject(root)
	if err != nil {
		return nil, false, fmt.Errorf("analyze project: %w", err)
	}
	if err := session.Save(root, sess); err != nil {
		return nil, false, fmt.Errorf("save session: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Session created: %d functions indexed (%s)\n\n",
		len(sess.Functions), sess.Lang)

	return sess, true, nil
}
