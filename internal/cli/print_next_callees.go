//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints all callees for the next function
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printNextCallees prints all callees for a function.
func printNextCallees(callees []model.Edge, sess *model.Session) {
	fmt.Printf("  callees (%d):\n", len(callees))
	for _, edge := range callees {
		printEdge("->", edge, sess)
	}
}
