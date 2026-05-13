//ff:func feature=endpoint type=implementation control=sequence
//ff:what Returns an error for unsupported languages
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

// Detect returns an error for unsupported languages.
func (d *UnsupportedDetector) Detect(_ string) ([]model.Endpoint, error) {
	return nil, &ErrUnsupported{Lang: d.Lang}
}
