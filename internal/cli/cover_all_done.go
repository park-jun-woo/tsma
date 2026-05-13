//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what Measures coverage for all DONE functions grouped by package
package cli

import (
	"fmt"
	"path/filepath"

	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/model"
)

// coverAllDone measures coverage for all DONE functions grouped by package.
func coverAllDone(root string, sess *model.Session, checker coverage.Checker) error {
	var doneFuncs []model.Function
	for _, fn := range sess.Functions {
		if fn.Status == model.StatusDone && fn.TestFile != "" {
			doneFuncs = append(doneFuncs, fn)
		}
	}

	if len(doneFuncs) == 0 {
		fmt.Println("No DONE functions with test files.")
		return nil
	}

	fmt.Printf("Measuring coverage for %d DONE functions...\n", len(doneFuncs))

	type pkgResult struct {
		pkgDir string
		count  int
		total  float64
	}
	pkgMap := map[string]*pkgResult{}
	var pkgOrder []string

	for i := range doneFuncs {
		fn := &doneFuncs[i]
		pkgDir := filepath.Dir(fn.File)

		if _, ok := pkgMap[pkgDir]; !ok {
			pkgMap[pkgDir] = &pkgResult{pkgDir: pkgDir}
			pkgOrder = append(pkgOrder, pkgDir)
		}

		report, err := checker.Check(root, fn.TestFile, fn)
		if err != nil {
			continue
		}

		pkgMap[pkgDir].count++
		pkgMap[pkgDir].total += report.TotalPct
	}

	var overallTotal float64
	var overallCount int

	for _, pkgDir := range pkgOrder {
		pr := pkgMap[pkgDir]
		if pr.count == 0 {
			continue
		}
		avg := pr.total / float64(pr.count)
		fmt.Printf("  %s/: %.0f%% (%d functions)\n", pr.pkgDir, avg, pr.count)
		overallTotal += pr.total
		overallCount += pr.count
	}

	if overallCount > 0 {
		fmt.Printf("\nOverall: %.0f%% average\n", overallTotal/float64(overallCount))
	}

	return nil
}
