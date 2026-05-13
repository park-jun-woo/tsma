//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Advances CurrentIndex past completed functions and returns the next TODO
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// advanceToNext skips PASS/DONE functions from CurrentIndex and returns the next TODO.
// Returns nil if all functions are complete.
func advanceToNext(sess *model.Session) *model.Function {
	for sess.CurrentIndex < len(sess.Functions) {
		fn := &sess.Functions[sess.CurrentIndex]
		if fn.Status == model.StatusTodo {
			return fn
		}
		sess.CurrentIndex++
	}
	return nil
}
