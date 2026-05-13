//ff:func feature=chain type=helper control=iteration dimension=1
//ff:what Checks if a receiver name indicates a repository/database boundary
package chain

import "strings"

// isRepoBoundary checks if the receiver name indicates a repository/database boundary.
func isRepoBoundary(receiver string) bool {
	lower := strings.ToLower(receiver)
	for _, indicator := range tsRepoIndicators {
		if strings.Contains(lower, indicator) {
			return true
		}
	}
	return false
}
