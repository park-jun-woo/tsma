//ff:type feature=index type=model
//ff:what Holds the mutable parsing state while indexing a Java file line by line
package index

import "github.com/park-jun-woo/tsma/internal/model"

// javaParseState carries the running state of the Java file indexer.
type javaParseState struct {
	relDir           string
	relPath          string
	pkg              string
	functions        []model.Function
	scopes           []javaScope
	depth            int
	lineNum          int
	lastNonEmptyLine int
}
