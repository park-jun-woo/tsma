package endpoint

import "regexp"

// Decorator patterns:
//   @app.get("/path")          @router.get("/path")
//   @app.post("/path")         @router.post("/path")
//   @app.api_route("/path", methods=["GET", "POST"])
var pyDecoratorRe = regexp.MustCompile(
	`^\s*@\s*\w+\.(get|post|put|patch|delete|head|options)\s*\(\s*["']([^"']+)["']`,
)

var pyAPIRouteRe = regexp.MustCompile(
	`^\s*@\s*\w+\.api_route\s*\(\s*["']([^"']+)["']\s*,\s*methods\s*=\s*\[([^\]]+)\]`,
)

var pyDefRe = regexp.MustCompile(`^(\s*)(?:async\s+)?def\s+(\w+)\s*\(`)

// pySkipDirs lists directories to skip when walking Python projects.
var pySkipDirs = map[string]bool{
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
