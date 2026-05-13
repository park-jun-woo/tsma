package endpoint

import "regexp"

// URL pattern regexes for urls.py parsing.
var (
	// path('route/', view_function, ...) or path('route/', view_function)
	djangoPathRe = regexp.MustCompile(
		`path\s*\(\s*["']([^"']*)["']\s*,\s*(\w[\w.]*)`,
	)
	// re_path(r'pattern', view_function, ...)
	djangoRePathRe = regexp.MustCompile(
		`re_path\s*\(\s*r?["']([^"']*)["']\s*,\s*(\w[\w.]*)`,
	)

	// def function_name(request
	djangoFuncViewRe = regexp.MustCompile(`^(\s*)def\s+(\w+)\s*\(\s*request`)

	// class ClassName(View / APIView / ...)
	djangoClassViewRe = regexp.MustCompile(`^(\s*)class\s+(\w+)\s*\(`)

	// HTTP method definitions inside class-based views
	djangoClassMethodRe = regexp.MustCompile(`^\s+def\s+(get|post|put|patch|delete|head|options)\s*\(\s*self`)
)

// djangoSkipDirs lists directories to skip when walking Django projects.
var djangoSkipDirs = map[string]bool{
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
