//ff:type feature=index type=model lang=csharp
//ff:what Holds the mutable parsing state while indexing a C# file line by line
package index

import "github.com/park-jun-woo/tsma/internal/model"

// csParseState carries the running state of the C# file indexer. A file-scoped
// namespace ("namespace X;") is recorded in fileNs; block-scoped namespaces and
// types are tracked on the scopes stack by brace depth.
type csParseState struct {
	relDir           string
	relPath          string
	fileNs           string
	functions        []model.Function
	scopes           []csScope
	pending          []string // namespace/type names awaiting their opening brace (Allman style)
	depth            int
	lineNum          int
	lastNonEmptyLine int
}
