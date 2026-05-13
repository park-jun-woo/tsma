//ff:func feature=endpoint type=implementation control=sequence
//ff:what Returns the error message for unsupported language detection
package endpoint

func (e *ErrUnsupported) Error() string {
	return "endpoint detection not implemented for: " + e.Lang
}
