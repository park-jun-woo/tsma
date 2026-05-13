//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints a page of functions with aligned name, status, callers, and callees
package cli

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// printFuncList prints a page of functions with aligned columns.
func printFuncList(functions []model.Function, maxName int) {
	for _, fn := range functions {
		status := strings.ToUpper(fn.Status)
		if fn.Dead {
			status = "DEAD"
		}
		callers := len(fn.Callers)
		callees := len(fn.Callees)
		fmt.Printf("  %-*s  %-7s  callers: %d  callees: %d\n",
			maxName, fn.Name, status, callers, callees)
	}
}
