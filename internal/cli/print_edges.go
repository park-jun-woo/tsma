//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints a list of call graph edges with direction label
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printEdges prints a labeled list of call graph edges.
func printEdges(label string, arrow string, edges []model.Edge, sess *model.Session) {
	fmt.Printf("  %s (%d):\n", label, len(edges))
	for _, edge := range edges {
		ambig := ""
		if edge.Ambiguous {
			ambig = " (?)"
		}
		target := sess.FindFunction(edge.Target)
		if target != nil {
			fmt.Printf("    %s %-30s %s:%d-%d%s\n",
				arrow, target.Name+"()", target.File, target.StartLine, target.EndLine, ambig)
		} else {
			fmt.Printf("    %s %s%s\n", arrow, edge.Target, ambig)
		}
	}
}
