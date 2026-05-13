//ff:func feature=runner type=implementation control=sequence
//ff:what Returns an error for languages without test runner support
package runner

import "fmt"

// Run returns an error for unsupported languages.
func (r *UnsupportedRunner) Run(_, _ string) (*Result, error) {
	return nil, fmt.Errorf("test runner not implemented for: %s", r.Lang)
}
