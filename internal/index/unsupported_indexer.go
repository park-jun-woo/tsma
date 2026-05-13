//ff:type feature=index type=implementation
//ff:what Returns an error for unsupported languages during indexing
package index

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// UnsupportedIndexer is returned for languages without indexer support.
type UnsupportedIndexer struct {
	Lang string
}

// Index returns an error indicating the language is not supported.
func (u *UnsupportedIndexer) Index(_ string) ([]model.Function, error) {
	return nil, fmt.Errorf("indexing not implemented for: %s", u.Lang)
}
