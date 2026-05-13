//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Returns a pointer to the first TODO or PARTIAL function in the session
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// findFirstIncomplete returns the first function with TODO or PARTIAL status.
func findFirstIncomplete(sess *model.Session) *model.Function {
	for i := range sess.Functions {
		if sess.Functions[i].Status == model.StatusTodo || sess.Functions[i].Status == model.StatusPartial {
			return &sess.Functions[i]
		}
	}
	return nil
}
