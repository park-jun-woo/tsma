//ff:func feature=graph type=helper control=sequence
//ff:what Defines regex patterns for extracting Python call expressions and imports
package graph

import "regexp"

var (
	// pyCallPattern matches function calls: identifier(, identifier.method(
	pyCallPattern = regexp.MustCompile(`(?:await\s+)?(\w+)(?:\.(\w+))?\s*\(`)
	// pyImportAbsPattern matches "import module" or "import module as alias".
	pyImportAbsPattern = regexp.MustCompile(`^import\s+(\w+)(?:\s+as\s+(\w+))?`)
	// pyFromImportPattern matches "from module import name" or "from .module import name".
	pyFromImportPattern = regexp.MustCompile(`^from\s+(\S+)\s+import\s+(.+)`)
)
