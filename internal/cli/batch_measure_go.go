//ff:func feature=cli type=helper control=iteration dimension=2
//ff:what Batch-measures all Go functions package-by-package, one go test run per package
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/tsma/internal/coverage"
	"github.com/park-jun-woo/tsma/internal/match"
	"github.com/park-jun-woo/tsma/internal/model"
	"github.com/park-jun-woo/tsma/internal/runner"
)

// batchMeasureGo measures every Go function by grouping them per package
// directory and running `go test -coverprofile` once per package. The -run
// filter is the union of the test functions attributed to that package's
// functions, so coverage is measured by the attributed tests (matching the
// single-function Check semantics, just batched). Untested functions are left
// TODO. A package whose test run fails (e.g. compile error) is skipped: its
// functions stay TODO and a failure is reported, without aborting the scan.
func batchMeasureGo(root string, sess *model.Session) {
	funcs := make([]*model.Function, len(sess.Functions))
	for i := range sess.Functions {
		funcs[i] = &sess.Functions[i]
	}
	batchMeasureGoFuncs(root, sess.Lang, funcs)
}

// batchMeasureGoFuncs is the core batch measurement over an explicit slice of
// function pointers, so it can measure either the whole session (first scan) or
// just the newly-reconciled subset (Phase012 rescan) without re-running tests
// for already-measured functions.
func batchMeasureGoFuncs(root, lang string, funcs []*model.Function) {
	matcher := match.NewFuncMatcher(lang)
	groups := map[string]*goPkgGroup{}
	var order []string

	for _, fn := range funcs {
		m, found := matcher.MatchFunc(root, fn)
		if !found || len(m.Files) == 0 {
			continue // untested: no test attributed, stays TODO (not measured)
		}
		pkgDir := goPkgDirOf(root, m)
		g := groups[pkgDir]
		if g == nil {
			g = &goPkgGroup{matches: map[*model.Function]match.TestMatch{}}
			groups[pkgDir] = g
			order = append(order, pkgDir)
		}
		g.funcs = append(g.funcs, fn)
		g.matches[fn] = m
		if tfs, err := runner.ResolveTestFuncs(m); err == nil {
			g.testFuncs = appendUnique(g.testFuncs, tfs)
		}
	}

	checker := &coverage.GoChecker{}
	total := len(order)
	for idx, pkgDir := range order {
		g := groups[pkgDir]
		fmt.Fprintf(os.Stderr, "Measuring coverage: pkg %d/%d (%s)\n", idx+1, total, shortPkg(root, pkgDir))
		reports, err := checker.CheckPkg(root, pkgDir, g.funcs, g.testFuncs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  skipped (%s): %s\n", shortPkg(root, pkgDir), firstLine(err.Error()))
			continue // leave this package's functions TODO
		}
		for _, fn := range g.funcs {
			report := reports[fn]
			if report == nil {
				continue
			}
			applyBatchReport(root, fn, g.matches[fn], report)
		}
	}
}
