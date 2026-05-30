//ff:func feature=cli type=helper control=sequence
//ff:what Appends the remaining-TODO summary when no TODO is measurable so the user can pick work
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// maybePrintNoProgressSummary prints the remaining-TODO summary when no TODO has
// a changed test, i.e. repeatedly calling `tsma next` can make no progress until
// the user edits something. It is appended after surfacing the current cursor
// TODO (the cursor still advances), so the rotation never traps a partial while
// the agent also sees the full list of work waiting. A no-op while progress is
// still possible.
func maybePrintNoProgressSummary(root string, sess *model.Session) {
	if anyTodoMeasurable(root, sess.Lang, sess) {
		return
	}
	printRemainingTodos(sess)
}
