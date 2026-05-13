//ff:func feature=chain type=implementation control=sequence
//ff:what Returns a human-readable error message for the unsupported language
package chain

// Error returns the error message.
func (e *ErrUnsupported) Error() string {
	return "chain tracing not implemented for: " + e.Lang
}
