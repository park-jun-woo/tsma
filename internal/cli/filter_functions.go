//ff:func feature=cli type=helper control=selection
//ff:what Filters functions by status using a switch on the status filter value
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// filterFunctions filters functions by status string.
func filterFunctions(functions []model.Function, statusFilter string) ([]model.Function, error) {
	switch statusFilter {
	case "dead":
		return filterByDead(functions), nil
	case "todo":
		return filterByStatus(functions, model.StatusTodo), nil
	case "partial":
		return filterByStatus(functions, model.StatusPartial), nil
	case "done":
		return filterByStatus(functions, model.StatusDone), nil
	case "":
		return filterNonDead(functions), nil
	default:
		return nil, fmt.Errorf("unknown status filter: %s (use: todo, partial, done, dead)", statusFilter)
	}
}
