//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Sets status and fail output for all functions in a group
package cli

import "github.com/park-jun-woo/tsma/internal/model"

// setGroupStatus sets the status and fail output for all functions in a group.
func setGroupStatus(funcs []*model.Function, status, failOutput string) {
	for _, fn := range funcs {
		fn.Status = status
		fn.FailOutput = failOutput
	}
}
