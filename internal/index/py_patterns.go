//ff:func feature=index type=helper control=sequence
//ff:what Defines regex patterns for matching Python function and class definitions
package index

import "regexp"

var (
	// pyFuncPattern matches def and async def declarations.
	pyFuncPattern = regexp.MustCompile(`^(\s*)(?:async\s+)?def\s+(\w+)\s*\(`)
	// pyClassDefPattern matches class declarations.
	pyClassDefPattern = regexp.MustCompile(`^(\s*)class\s+(\w+)`)
)
