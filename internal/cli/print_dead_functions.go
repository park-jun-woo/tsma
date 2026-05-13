//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints dead code functions from the session
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printDeadFunctions prints all dead code functions.
func printDeadFunctions(functions []model.Function) {
	count := 0
	for _, fn := range functions {
		if !fn.Dead {
			continue
		}
		count++
		fmt.Printf("  %-30s %s:%d-%d\n", fn.Name, fn.File, fn.StartLine, fn.EndLine)
	}
	if count == 0 {
		fmt.Println("No dead code found.")
	} else {
		fmt.Printf("\n%d dead function(s)\n", count)
	}
}
