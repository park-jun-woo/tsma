//ff:func feature=cli type=helper control=sequence
//ff:what Measures and prints coverage for a single named function
package cli

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/model"
)

// coverSingleFunction measures and prints coverage for a single function.
func coverSingleFunction(root string, sess *model.Session, checker coverage.Checker, name string) error {
	fn := sess.FindFunction(name)
	if fn == nil {
		return fmt.Errorf("function not found: %s", name)
	}
	if fn.TestFile == "" {
		return fmt.Errorf("no test file for function: %s", name)
	}

	report, err := checker.Check(root, fn.TestFile, fn)
	if err != nil {
		return fmt.Errorf("check coverage: %w", err)
	}

	fmt.Printf("%s (%s:%d-%d)\n", fn.Name, fn.File, fn.StartLine, fn.EndLine)

	if len(report.Funcs) > 0 {
		fc := report.Funcs[0]
		fmt.Printf("  coverage: %.0f%% (%d/%d branches)\n",
			fc.CoveredPct, fc.CoveredBlocks, fc.TotalBlocks)
		for _, line := range fc.UncoveredLines {
			fmt.Printf("    UNCOVERED: line %d\n", line)
		}
	} else {
		fmt.Printf("  coverage: %.0f%%\n", report.TotalPct)
		for _, ub := range report.Uncovered {
			fmt.Printf("    UNCOVERED: line %d\n", ub.Line)
		}
	}

	return nil
}
