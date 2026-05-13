//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Classifies a TS/JS call receiver as repository-interface or external based on naming heuristics
package chain

import "strings"

// classifyTSBoundary determines the boundary type for an external call based on the receiver name.
func classifyTSBoundary(receiver string) string {
	lower := strings.ToLower(receiver)
	for _, indicator := range tsRepoIndicators {
		if strings.Contains(lower, indicator) {
			return "repository-interface"
		}
	}
	return "external"
}
