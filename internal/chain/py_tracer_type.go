//ff:type feature=chain type=model
//ff:what Regex-based Python call chain tracer with associated constants and patterns
package chain

import "regexp"

// PyTracer traces function call chains in Python projects using regex-based analysis.
type PyTracer struct{}

// maxPyTraceDepth is the recursion limit for call chain tracing.
const maxPyTraceDepth = 10

// pyTracerSkipDirs lists directories to skip when indexing Python files.
var pyTracerSkipDirs = map[string]bool{
	"__pycache__":   true,
	".venv":         true,
	"venv":          true,
	".git":          true,
	".tsma":   true,
	"node_modules":  true,
	"site-packages": true,
	"dist":          true,
	"build":         true,
}

// Regex patterns for Python call detection.
var (
	// Matches: identifier(, identifier.method(, await identifier(, await identifier.method(
	pyCallRe = regexp.MustCompile(`(?:await\s+)?(\w+(?:\.\w+)*)\s*\(`)

	// Matches: def funcName( or async def funcName(
	pyFuncDefRe = regexp.MustCompile(`^(\s*)(?:async\s+)?def\s+(\w+)\s*\(`)

	// Matches: class ClassName(... or class ClassName:
	pyClassDefRe = regexp.MustCompile(`^(\s*)class\s+(\w+)\s*[\(:]`)

	// Matches import patterns to determine if something is external.
	// from .module import ... (relative import)
	pyRelativeImportRe = regexp.MustCompile(`^\s*from\s+\.`)
)
