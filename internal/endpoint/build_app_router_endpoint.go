//ff:func feature=endpoint type=helper control=sequence
//ff:what Builds an endpoint from an App Router regex match location
package endpoint

import "github.com/park-jun-woo/tsma/internal/model"

func buildAppRouterEndpoint(content string, lines []string, loc []int, relPath, routePath, baseName string) model.Endpoint {
	method := content[loc[2]:loc[3]]
	lineNum := countNewlines(content[:loc[0]]) + 1
	startLine, endLine := findExportedFuncBounds(lines, lineNum-1)

	return model.Endpoint{
		Name:   baseName + method,
		Method: method,
		Path:   routePath,
		Handler: model.FuncLocation{
			File:      relPath,
			StartLine: startLine,
			EndLine:   endLine,
		},
		Status: model.StatusTodo,
	}
}
