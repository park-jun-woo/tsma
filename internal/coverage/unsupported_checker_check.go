//ff:func feature=coverage type=implementation control=sequence
//ff:what Returns an ErrUnsupported error for languages without coverage support
package coverage

import "github.com/park-jun-woo/tsma/internal/model"

// Check returns an error for unsupported languages.
func (c *UnsupportedChecker) Check(_, _ string, _ *model.Function) (*Report, error) {
	return nil, &ErrUnsupported{Lang: c.Lang}
}
