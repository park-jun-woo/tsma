//ff:type feature=endpoint type=model
//ff:what Stores a parsed Python decorator waiting for its def line
package endpoint

// pendingDecorator stores a parsed decorator waiting for its def line.
type pendingDecorator struct {
	method string
	path   string
	line   int
}
