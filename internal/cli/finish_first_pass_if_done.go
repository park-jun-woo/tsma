//ff:func feature=cli type=helper control=sequence
//ff:what Flips the session into interactive mode once the first-pass watermark reaches the end
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// finishFirstPassIfDone marks the first pass complete and resets the cursor to 0
// once CurrentIndex (the first-pass watermark) has advanced past every function.
// After this, runNext switches from the forward-only watermark to the rotating
// interactive cursor. No-op if the first pass is already done or still running.
func finishFirstPassIfDone(sess *model.Session) {
	if sess.FirstPassDone {
		return
	}
	if sess.CurrentIndex >= len(sess.Functions) {
		sess.FirstPassDone = true
		sess.CurrentIndex = 0
	}
}
