//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints up to a limited number of callers for the next function
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printNextCallers prints up to limit callers for a function.
func printNextCallers(callers []model.Edge, limit int, sess *model.Session) {
	shown := len(callers)
	if shown > limit {
		shown = limit
	}
	fmt.Printf("  callers (%d shown):\n", shown)
	for i := 0; i < shown; i++ {
		printEdge("<-", callers[i], sess)
	}
}
