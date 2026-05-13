//ff:type feature=endpoint type=model
//ff:what Error type indicating an unsupported language for endpoint detection
package endpoint

// ErrUnsupported indicates an unsupported language.
type ErrUnsupported struct {
	Lang string
}
