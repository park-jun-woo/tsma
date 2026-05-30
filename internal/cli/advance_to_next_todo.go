//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Rotating-cursor TODO finder used after the first pass: scans from CurrentIndex, wrapping
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// advanceToNextTodo returns the next StatusTodo function starting at CurrentIndex,
// wrapping around to index 0 once. It moves CurrentIndex onto the returned
// function's index (the rotating cursor). Returns nil if no TODO exists.
//
// Unlike advanceToNext (the first-pass watermark, which only moves forward and
// stops at the end), this is the interactive-mode cursor: it never gets stuck at
// the end, so a single unchanged partial can never gate the TODOs behind it.
func advanceToNextTodo(sess *model.Session) *model.Function {
	n := len(sess.Functions)
	if n == 0 {
		return nil
	}
	if sess.CurrentIndex < 0 || sess.CurrentIndex >= n {
		sess.CurrentIndex = 0
	}
	for off := 0; off < n; off++ {
		i := (sess.CurrentIndex + off) % n
		if sess.Functions[i].Status == model.StatusTodo {
			sess.CurrentIndex = i
			return &sess.Functions[i]
		}
	}
	return nil
}
