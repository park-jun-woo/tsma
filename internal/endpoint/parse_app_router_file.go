//ff:func feature=endpoint type=implementation control=iteration dimension=1
//ff:what Reads a route.ts file and finds exported GET/POST/etc functions with path derivation
package endpoint

import (
	"os"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

// parseAppRouterFile parses an App Router route file for exported HTTP method handlers.
func parseAppRouterFile(filePath, relPath string) ([]model.Endpoint, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	routePath := deriveRoutePath(relPath)
	baseName := deriveEndpointName(relPath)

	var endpoints []model.Endpoint

	matches := nextjsExportPattern.FindAllStringSubmatchIndex(content, -1)
	for _, loc := range matches {
		ep := buildAppRouterEndpoint(content, lines, loc, relPath, routePath, baseName)
		endpoints = append(endpoints, ep)
	}

	return endpoints, nil
}
