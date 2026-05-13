//ff:func feature=cli type=helper control=selection
//ff:what Runs the test for a package group and updates function statuses accordingly
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/runner"
)

// applyCheckResult runs the test for a package group and updates function statuses.
func applyCheckResult(root string, r runner.Runner, g pkgGroup) {
	if g.testFile == "" {
		fmt.Fprintf(os.Stderr, "  %s/: no test files → TODO\n", g.pkgDir)
		return
	}

	result, err := r.Run(root, g.testFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s/: run error: %v\n", g.pkgDir, err)
		setGroupStatus(g.funcs, model.StatusFail, err.Error())
		return
	}

	switch {
	case result.Pass:
		fmt.Fprintf(os.Stderr, "  %s/: %d functions, PASS → DONE\n", g.pkgDir, len(g.funcs))
		setGroupStatus(g.funcs, model.StatusDone, "")
	default:
		fmt.Fprintf(os.Stderr, "  %s/: %d functions, FAIL\n", g.pkgDir, len(g.funcs))
		printFailSummary(result.Output)
		setGroupStatus(g.funcs, model.StatusFail, result.Output)
	}
}
