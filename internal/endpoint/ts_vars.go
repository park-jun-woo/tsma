package endpoint

import "regexp"

// expressRoutePattern matches Express route registration patterns like:
// app.get("/path", handler), router.post("/path", handler)
var expressRoutePattern = regexp.MustCompile(
	`(?:app|router)\.(get|post|put|patch|delete|all)\(\s*["']([^"']+)["']\s*,\s*([a-zA-Z_$][\w.$]*)`,
)

// tsFuncDefPattern matches TypeScript/JavaScript function definitions.
var tsFuncDefPattern = regexp.MustCompile(
	`(?:export\s+)?(?:async\s+)?function\s+(\w+)|(?:export\s+)?(?:const|let|var)\s+(\w+)\s*=`,
)

// tsSkipDirs lists directories to skip when walking TypeScript projects.
var tsSkipDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	".tsma":  true,
	"dist":         true,
	"build":        true,
}
