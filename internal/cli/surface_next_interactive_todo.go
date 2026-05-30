//ff:func feature=cli type=helper control=sequence
//ff:what After a PASS/DONE, surfaces the next interactive TODO via the rotating cursor or the all-complete banner
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// surfaceNextInteractiveTodo advances the rotating cursor to the next TODO and
// prints it as the continuation target. If none remain it prints the
// all-complete banner. Used by the interactive PASS/DONE branches.
func surfaceNextInteractiveTodo(sess *model.Session) {
	if next := advanceToNextTodo(sess); next != nil {
		printContinueInstruction()
		printTodoFunction(next, next.TestFile)
		printNextInstruction()
	} else {
		printAllComplete()
	}
}
