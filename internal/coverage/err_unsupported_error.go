//ff:func feature=coverage type=implementation control=sequence
//ff:what Returns a human-readable error message for the unsupported language
package coverage

// Error returns the error message.
func (e *ErrUnsupported) Error() string {
	return "coverage check not implemented for: " + e.Lang
}
