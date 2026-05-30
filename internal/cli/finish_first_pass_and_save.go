//ff:func feature=cli type=helper control=sequence
//ff:what Finalizes a non-measuring first-pass step: flip-if-done, recount, persist
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// finishFirstPassAndSave wraps the tail shared by the first pass's non-measuring
// branches (no-test, unchanged): reconcile the first-pass flag now that the
// watermark advanced, recompute the summary, and persist.
func finishFirstPassAndSave(root string, sess *model.Session) error {
	finishFirstPassIfDone(sess)
	sess.RecalcSummary()
	return saveSession(root, sess)
}
