//ff:func feature=cli type=helper control=selection
//ff:what Sorts functions by the given sort field using a switch dispatch
package cli

import (
	"fmt"
	"sort"

	"github.com/park-jun-woo/tsma/internal/model"
)

// sortFunctions sorts the function slice in place by the given field.
func sortFunctions(functions []model.Function, sortField string) error {
	switch sortField {
	case "priority":
		sort.Slice(functions, func(i, j int) bool {
			ci := len(functions[i].Callers)
			cj := len(functions[j].Callers)
			if ci != cj {
				return ci > cj
			}
			return functions[i].Name < functions[j].Name
		})
	case "name":
		sort.Slice(functions, func(i, j int) bool {
			return functions[i].Name < functions[j].Name
		})
	case "file":
		sort.Slice(functions, func(i, j int) bool {
			if functions[i].File != functions[j].File {
				return functions[i].File < functions[j].File
			}
			return functions[i].StartLine < functions[j].StartLine
		})
	default:
		return fmt.Errorf("unknown sort field: %s (use: priority, name, file)", sortField)
	}
	return nil
}
