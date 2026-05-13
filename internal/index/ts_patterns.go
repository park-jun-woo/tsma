//ff:func feature=index type=helper control=sequence
//ff:what Defines regex patterns for matching TS/JS function and class declarations
package index

import "regexp"

var (
	// tsFuncPattern matches function declarations and const arrow functions.
	tsFuncPattern = regexp.MustCompile(
		`^(?:export\s+)?(?:async\s+)?function\s+(\w+)\s*\(` +
			`|^(?:export\s+)?(?:const|let|var)\s+(\w+)\s*=`,
	)
	// tsClassPattern matches class declarations.
	tsClassPattern = regexp.MustCompile(`^(?:export\s+)?class\s+(\w+)`)
	// tsMethodPattern matches class method declarations (indented).
	tsMethodPattern = regexp.MustCompile(`^\s+(?:async\s+)?(\w+)\s*\([^)]*\)\s*[:{]`)
)
