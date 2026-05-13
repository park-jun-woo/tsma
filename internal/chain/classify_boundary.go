//ff:func feature=chain type=helper control=selection
//ff:what Determines boundary type for an external Go call based on receiver name heuristics
package chain

import "strings"

// classifyBoundary determines the boundary type for an external call.
func classifyBoundary(receiver string) string {
	lower := strings.ToLower(receiver)
	switch {
	case strings.Contains(lower, "repo"):
		return "repository-interface"
	case strings.Contains(lower, "db"):
		return "database"
	case strings.Contains(lower, "store"):
		return "repository-interface"
	default:
		return "external"
	}
}
