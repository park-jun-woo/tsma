//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Reports whether any TODO function has a changed test that could be measured this lap
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// anyTodoMeasurable returns true if at least one StatusTodo function has a test
// whose change detectTestChange would report (changed=true). It is used in
// interactive mode to decide whether a full rotation can make progress: when no
// TODO is measurable, repeatedly calling `tsma next` would only rotate the
// cursor forever, so the caller instead prints a remaining-TODO summary and
// terminates normally.
func anyTodoMeasurable(root, lang string, sess *model.Session) bool {
	for i := range sess.Functions {
		fn := &sess.Functions[i]
		if fn.Status != model.StatusTodo {
			continue
		}
		if changed, _ := detectTestChange(root, lang, fn); changed {
			return true
		}
	}
	return false
}
