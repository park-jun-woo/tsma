//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Iterates over package groups and applies check results for each
package cli

import "github.com/park-jun-woo/tsma/internal/runner"

// checkAllGroups iterates over package groups and applies check results.
func checkAllGroups(root string, r runner.Runner, groups []pkgGroup) {
	for _, g := range groups {
		applyCheckResult(root, r, g)
	}
}
