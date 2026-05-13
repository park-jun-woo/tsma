//ff:type feature=chain type=model
//ff:what Regex-based TypeScript/JavaScript call chain tracer with associated patterns and constants
package chain

import "regexp"

// TSTracer traces function call chains in TypeScript/JavaScript using regex-based analysis.
type TSTracer struct{}

// tsFuncCallPattern matches common function call patterns in TypeScript/JavaScript.
// Captures: optional receiver (group 1), function/method name (group 2).
var tsFuncCallPattern = regexp.MustCompile(
	`(?:await\s+)?([a-zA-Z_$][\w$]*)\s*\.\s*([a-zA-Z_$][\w$]*)\s*\(` +
		`|(?:await\s+)?([a-zA-Z_$][\w$]*)\s*\(`,
)

// tsSkipCallNames lists identifiers that should not be traced (language builtins, control flow).
var tsSkipCallNames = map[string]bool{
	"console": true, "require": true, "import": true,
	"JSON": true, "Math": true, "Date": true, "Array": true,
	"Object": true, "String": true, "Number": true, "Boolean": true,
	"Promise": true, "Error": true, "RegExp": true, "Map": true,
	"Set": true, "parseInt": true, "parseFloat": true,
	"setTimeout": true, "setInterval": true, "clearTimeout": true,
	"clearInterval": true, "fetch": true, "Buffer": true,
	"if": true, "else": true, "for": true, "while": true,
	"switch": true, "return": true, "throw": true, "new": true,
}

// tsRepoIndicators lists substrings in identifiers that suggest a repository/database boundary.
var tsRepoIndicators = []string{"repo", "db", "store", "prisma", "knex", "sequelize", "mongoose", "typeorm"}

// maxTraceDepth limits recursion depth to prevent excessive tracing.
const maxTraceDepth = 10
