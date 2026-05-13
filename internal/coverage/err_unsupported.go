//ff:type feature=coverage type=model
//ff:what Error type indicating a language is not supported for coverage checking
package coverage

// ErrUnsupported indicates an unsupported language.
type ErrUnsupported struct {
	Lang string
}
