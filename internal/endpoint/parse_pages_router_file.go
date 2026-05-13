//ff:func feature=endpoint type=implementation control=sequence
//ff:what Reads a pages/api file and finds export default to derive path and name
package endpoint

import (
	"os"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// parsePagesRouterFile parses a Pages Router API file for a default export handler.
func parsePagesRouterFile(filePath, relPath string) (*model.Endpoint, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	loc := nextjsDefaultExportPattern.FindStringSubmatchIndex(content)
	if loc == nil {
		return nil, nil
	}

	lineNum := countNewlines(content[:loc[0]]) + 1
	startLine, endLine := findExportedFuncBounds(lines, lineNum-1)

	routePath := deriveRoutePath(relPath)
	epName := deriveEndpointName(relPath)

	return &model.Endpoint{
		Name:   epName,
		Method: "ANY",
		Path:   routePath,
		Handler: model.FuncLocation{
			File:      relPath,
			StartLine: startLine,
			EndLine:   endLine,
		},
		Status: model.StatusTodo,
	}, nil
}
