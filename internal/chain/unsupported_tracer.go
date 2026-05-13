//ff:type feature=chain type=model
//ff:what Placeholder tracer for languages without chain tracing support
package chain

// UnsupportedTracer is a placeholder for unimplemented languages.
type UnsupportedTracer struct {
	Lang string
}
