//ff:type feature=index type=factory
//ff:what Defines the Indexer interface for language-specific function indexing
package index

import "github.com/park-jun-woo/tsma/internal/model"

// Indexer collects all function declarations from a project.
type Indexer interface {
	Index(projectRoot string) ([]model.Function, error)
}
