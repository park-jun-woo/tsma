//ff:func feature=cli type=helper control=sequence
//ff:what Prints callers and callees for a specific function in the call graph
package cli

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printFuncGraph prints the graph edges for a specific function.
func printFuncGraph(sess *model.Session, fn *model.Function) {
	status := strings.ToUpper(fn.Status)
	if fn.Dead {
		status = "DEAD"
	}
	fmt.Printf("%s (%s:%d-%d) — %s\n", fn.Name, fn.File, fn.StartLine, fn.EndLine, status)

	showCallers := !graphCallees || graphCallers
	showCallees := !graphCallers || graphCallees

	if showCallers && len(fn.Callers) > 0 {
		printEdges("callers", "<-", fn.Callers, sess)
	} else if showCallers {
		fmt.Println("  callers: (none)")
	}

	if showCallees && len(fn.Callees) > 0 {
		printEdges("callees", "->", fn.Callees, sess)
	} else if showCallees {
		fmt.Println("  callees: (none)")
	}
}
