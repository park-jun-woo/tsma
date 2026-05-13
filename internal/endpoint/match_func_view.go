//ff:func feature=endpoint type=implementation control=sequence
//ff:what Matches a Django function-based view definition and returns endpoint
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

func matchFuncView(line string, lines []string, idx, lineNum int, relPath string, route djangoRoute) []model.Endpoint {
	m := djangoFuncViewRe.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	if m[2] != route.viewName {
		return nil
	}

	endLine := findPyFuncEndDjango(lines, idx, m[1])
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
