//ff:func feature=endpoint type=implementation control=iteration dimension=1
//ff:what Regex-scans a TS/JS file for Express route registrations and resolves handler locations
package endpoint

import (
	"os"
	"path/filepath"
	"strings"
)

// parseExpressRoutes parses a single file for Express route registrations.
func parseExpressRoutes(filePath, projectRoot string) ([]tsRouteRegistration, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	relPath, _ := filepath.Rel(projectRoot, filePath)

	var regs []tsRouteRegistration

	matches := expressRoutePattern.FindAllStringSubmatchIndex(content, -1)
	for _, loc := range matches {
		reg := buildExpressRegistration(content, lines, loc, relPath)
		regs = append(regs, reg)
	}

	return regs, nil
}
