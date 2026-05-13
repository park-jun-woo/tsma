//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints a page of functions with aligned name and status columns
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
		fmt.Printf("  %-*s  %s\n", maxName, fn.Name, status)
	}
}
