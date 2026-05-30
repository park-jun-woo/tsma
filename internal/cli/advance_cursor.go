//ff:func feature=cli type=helper control=sequence
//ff:what Moves the interactive rotating cursor forward one slot, wrapping to 0 at the end
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// advanceCursor moves CurrentIndex forward one position, wrapping to 0 past the
// end. Used in interactive mode after surfacing a non-measuring TODO (or a
// PASS/DONE) so the next call inspects the next function rather than re-pinning
// the current one.
func advanceCursor(sess *model.Session) {
	n := len(sess.Functions)
	if n == 0 {
		return
	}
	sess.CurrentIndex = (sess.CurrentIndex + 1) % n
}
