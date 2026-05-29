//ff:type feature=index type=model
//ff:what Holds the mutable parsing state while indexing a Rust file line by line
package index

import "github.com/park-jun-woo/tsma/internal/model"

// rsParseState carries the running state of the Rust file indexer.
type rsParseState struct {
	relDir           string
	relPath          string
	functions        []model.Function
	scopes           []rsScope
	depth            int
	lineNum          int
	lastNonEmptyLine int
	pendingCfgTest   bool
}
