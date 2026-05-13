//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Prints detailed chain-level branch coverage for a specific endpoint
package cli

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/tsma/internal/model"
)

func showEndpointStatus(sess *model.Session, name string) error {
	ep := sess.FindEndpoint(name)
	if ep == nil {
		return fmt.Errorf("endpoint not found: %s", name)
	}

	status := strings.ToUpper(ep.Status)
	fmt.Printf("%s (%s:%d-%d) — %s\n",
		ep.Name, ep.Handler.File, ep.Handler.StartLine, ep.Handler.EndLine, status)

	if ep.TestFile != "" {
		fmt.Printf("  test file: %s\n", ep.TestFile)
	}

	for key, pct := range ep.Coverage {
		fmt.Printf("  chain coverage: %-40s %.0f%%\n", key, pct)
	}

	for _, line := range ep.UncoveredBranches {
		fmt.Printf("  uncovered branch: line %d\n", line)
	}

	if len(ep.Chain) == 0 {
		return nil
	}

	fmt.Println("  chain:")
	for _, ce := range ep.Chain {
		if ce.File != "" {
			fmt.Printf("    -> %-30s %s:%d-%d\n", ce.Func+"()", ce.File, ce.StartLine, ce.EndLine)
			continue
		}
		if ce.Boundary != "" {
			fmt.Printf("    -> %-30s (%s)\n", ce.Func+"()", ce.Boundary)
		}
	}

	return nil
}
