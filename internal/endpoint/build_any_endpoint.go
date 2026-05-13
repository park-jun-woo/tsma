//ff:func feature=endpoint type=helper control=sequence
//ff:what Builds a single ANY-method endpoint for a class view without specific HTTP methods
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

func buildAnyEndpoint(lines []string, idx, lineNum int, classIndent, relPath string, route djangoRoute) []model.Endpoint {
	endLine := findPyFuncEndDjango(lines, idx, classIndent)
	return []model.Endpoint{{
		Name:   route.viewName,
		Method: "ANY",
		Path:   route.path,
		Handler: model.FuncLocation{
			File:      relPath,
			StartLine: lineNum,
			EndLine:   endLine,
		},
		Status: model.StatusTodo,
	}}
}
