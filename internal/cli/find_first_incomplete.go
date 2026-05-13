//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Finds the first TODO or FAIL function in the session
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// findFirstIncomplete returns the first TODO or FAIL function, or nil if all DONE.
func findFirstIncomplete(sess *model.Session) *model.Function {
	for i := range sess.Functions {
		fn := &sess.Functions[i]
		if fn.Status == model.StatusTodo || fn.Status == model.StatusFail {
			return fn
		}
	}
	return nil
}
