//ff:type feature=chain type=model
//ff:what Error type indicating a language is not supported for chain tracing
package chain

// ErrUnsupported indicates an unsupported language.
type ErrUnsupported struct {
	Lang string
}
