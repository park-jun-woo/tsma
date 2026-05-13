//ff:func feature=chain type=implementation control=sequence
//ff:what Returns an ErrUnsupported error for languages without chain tracing
package chain

import "github.com/park-jun-woo/tsma/internal/model"

// Trace returns an error for unsupported languages.
func (t *UnsupportedTracer) Trace(_ string, _ model.FuncLocation) ([]model.ChainEntry, error) {
	return nil, &ErrUnsupported{Lang: t.Lang}
}
