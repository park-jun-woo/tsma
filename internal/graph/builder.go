//ff:type feature=graph type=factory
//ff:what Defines the Builder interface for language-specific call graph construction
package graph

import "github.com/park-jun-woo/tsma/internal/model"

// Builder constructs a call graph from indexed functions.
type Builder interface {
	Build(projectRoot string, functions []model.Function) ([]model.Function, model.GraphSummary, error)
}
