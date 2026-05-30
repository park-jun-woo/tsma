//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints a summary of the remaining TODO functions when interactive rotation can make no progress
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printRemainingTodos prints a summary of every TODO function left in the
// session. It is shown after the first pass when a full rotation can make no
// progress (no TODO has a changed test): the user/agent chooses which function
// to work on. Partials list their coverage; functions with no test are flagged
// as needing one. This is a normal terminal state ("waiting on the user"), not
// an error — complete (TODO 0) remains the only completion invariant.
func printRemainingTodos(sess *model.Session) {
	count := 0
	for i := range sess.Functions {
		if sess.Functions[i].Status == model.StatusTodo {
			count++
		}
	}
	fmt.Printf("\n%d TODO function(s) remaining — pick one to work on:\n", count)
	for i := range sess.Functions {
		fn := &sess.Functions[i]
		if fn.Status != model.StatusTodo {
			continue
		}
		if fn.CoveragePct > 0 {
			fmt.Printf("  %s  %.0f%% (improve coverage)  %s:%d-%d\n",
				fn.Name, fn.CoveragePct, fn.File, fn.StartLine, fn.EndLine)
		} else {
			fmt.Printf("  %s  (write test)  %s:%d-%d\n",
				fn.Name, fn.File, fn.StartLine, fn.EndLine)
		}
	}
	fmt.Println("  ▶ Improve a test or write a missing one, then run `tsma next`.")
}
