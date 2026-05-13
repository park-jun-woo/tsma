//ff:func feature=graph type=helper control=sequence
//ff:what Defines regex patterns for extracting TS/JS call expressions and imports
package graph

import "regexp"

var (
	// tsCallPattern matches function calls: identifier(, identifier.method(, await identifier(
	tsCallPattern = regexp.MustCompile(`(?:await\s+)?(\w+)(?:\.(\w+))?\s*\(`)
	// tsNamedImportPattern extracts individual imported names.
	tsNamedImportPattern = regexp.MustCompile(`import\s+\{([^}]*)\}\s+from\s+['"]([^'"]+)['"]`)
)
