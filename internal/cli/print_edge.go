//ff:func feature=cli type=helper control=sequence
//ff:what Prints a single call graph edge with direction arrow
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printEdge prints a single edge with direction arrow.
func printEdge(arrow string, edge model.Edge, sess *model.Session) {
	ambig := ""
	if edge.Ambiguous {
		ambig = " (?)"
	}
	target := sess.FindFunction(edge.Target)
	if target != nil {
		fmt.Printf("    %s %-20s %s:%d-%d%s\n",
			arrow, target.Name+"()", target.File, target.StartLine, target.EndLine, ambig)
	} else {
		fmt.Printf("    %s %s%s\n", arrow, edge.Target, ambig)
	}
}
