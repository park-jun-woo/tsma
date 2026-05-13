//ff:type feature=graph type=implementation
//ff:what Returns an error for unsupported languages during graph building
package graph

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// UnsupportedBuilder is returned for languages without graph builder support.
type UnsupportedBuilder struct {
	Lang string
}

// Build returns an error indicating the language is not supported.
func (u *UnsupportedBuilder) Build(_ string, _ []model.Function) ([]model.Function, model.GraphSummary, error) {
	return nil, model.GraphSummary{}, fmt.Errorf("graph building not implemented for: %s", u.Lang)
}
