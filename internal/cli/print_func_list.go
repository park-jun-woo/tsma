//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints a page of functions with aligned name, status, and coverage
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
		covStr := ""
		if fn.Status == model.StatusDone || fn.Status == model.StatusPartial {
			covStr = fmt.Sprintf("  %.0f%%", fn.CoveragePct)
		}
		fmt.Printf("  %-*s  %-7s%s\n", maxName, fn.Name, status, covStr)
	}
}
