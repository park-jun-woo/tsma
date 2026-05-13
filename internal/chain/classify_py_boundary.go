//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Determines boundary type for an external Python call based on naming heuristics
package chain

import "strings"

// classifyPyBoundary determines the boundary type for an external Python call.
func classifyPyBoundary(callExpr string, imports map[string]string) string {
	lower := strings.ToLower(callExpr)

	// Check repository-related patterns.
	repoPatterns := []string{"repo", "db", "store", "model", "orm", "session", "query"}
	for _, pattern := range repoPatterns {
		if strings.Contains(lower, pattern) {
			return "repository-interface"
		}
	}

	// Check if the root identifier is from an external import.
	parts := strings.Split(callExpr, ".")
	if len(parts) > 0 {
		root := parts[0]
		if src, ok := imports[root]; ok && src == "external" {
			return "external"
		}
	}

	return "external"
}
