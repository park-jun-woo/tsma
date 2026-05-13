//ff:type feature=chain type=model
//ff:what Defines the Tracer interface for language-specific call chain tracing
package chain

import "github.com/park-jun-woo/tsma/internal/model"

// Tracer traces the call chain from a handler function.
type Tracer interface {
	Trace(projectRoot string, handler model.FuncLocation) ([]model.ChainEntry, error)
}
